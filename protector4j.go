package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type Task struct {
	Email string `yaml:"email,omitempty"`
	Password string `yaml:"password,omitempty"`
	JarsPath []string `yaml:"jarsPath,omitempty"`
	ProtectAll bool `yaml:"protectAll,omitempty"`
	ProtectInnerJars bool `yaml:"protectInnerJars,omitempty"`
	ClassesToProtect []string `yaml:"classesToProtect,omitempty"`
	Exclude []string `yaml:"exclude,omitempty"`
	OutputFolder string `yaml:"outputFolder,omitempty"`
	TempFolder string `yaml:"tempFolder,omitempty"`
	JavaVersion string `yaml:"javaVersion,omitempty"`
	IncludeJavaFX bool `yaml:"includeJavaFX,omitempty"`
	KeySeed string `yaml:"keySeed,omitempty"`
	TargetPlatforms []string `yaml:"targetPlatforms,omitempty"`
}

func main() {
	app := &cli.App{
		Name: "vlinx-protector4j",
		Description: "Protect Java App from Decompilation, beyond Obfuscation",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "version",
				Required: true,
				Usage: "jre-version",
			},
			&cli.StringFlag{
				Name: "email",
				Required: true,
				Usage: "account-email",
			},
			&cli.StringFlag{
				Name: "password",
				Required: true,
				Usage: "md5-of-password",
			},
			&cli.BoolFlag{
				Name: "protect-all",
			},
			&cli.BoolFlag{
				Name: "protect-inner-jars",
			},
			&cli.StringFlag{
				Name: "classes-to-protect",
				Usage: "vlinx.test.TestClass1,vlinx.test.pack1.*,vlinx.test.pack1.**",
			},
			&cli.StringFlag{
				Name: "exclude",
				Usage: "vlinx.test.TestClass1,vlinx.test.pack1.*,vlinx.test.pack1.**",
			},
			&cli.StringFlag{
				Name: "output-folder",
				Value: ".",
			},
			&cli.BoolFlag{
				Name: "include-java-fx",
			},
			&cli.StringFlag{
				Name: "key-seed",
			},
			&cli.StringFlag{
				Name: "target-platforms",
				Usage: "linux64,win64,mac,linux32,win32",
			},
		},
		Action: func(context *cli.Context) error {
			task := new(Task)

			task.Email = context.String("email")
			task.Password = context.String("password")
			if task.Password == "-" {
				fmt.Scanf("%s", &task.Password)
			}

			for _, arg := range context.Args().Slice() {
				jarPath, err := filepath.Abs(arg)
				if err != nil {
					return err
				}

				task.JarsPath = append(task.JarsPath, jarPath)
			}

			task.ProtectAll = context.Bool("protect-all")
			task.ProtectInnerJars = context.Bool("protect-inner-jars")
			if context.IsSet("classes-to-protect") {
				task.ClassesToProtect = strings.Split(context.String("classes-to-protect"), ",")
			}
			if context.IsSet("exclude") {
				task.Exclude = strings.Split(context.String("exclude"), ",")
			}

			output, err := filepath.Abs(context.String("output-folder"))
			if err != nil {
				return err
			}
			temp, err := ioutil.TempDir(output, "vlinx-")
			if err != nil {
				return err
			}
			defer os.RemoveAll(temp)

			task.OutputFolder = temp
			task.TempFolder = temp

			task.JavaVersion = context.String("version")
			if strings.HasPrefix(task.JavaVersion, "8") {
				task.JavaVersion = "java-8"
			}
			if strings.HasPrefix(task.JavaVersion, "11") {
				task.JavaVersion = "java-11"
			}

			task.IncludeJavaFX = context.Bool("include-java-fx")

			task.KeySeed = context.String("key-seed")

			if context.IsSet("target-platforms") {
				task.TargetPlatforms = strings.Split(context.String("target-platforms"), ",")
			}

			bytes, err := yaml.Marshal(task)
			if err != nil {
				return err
			}
			if err := ioutil.WriteFile(filepath.Join(temp, "task.java.yaml"), bytes, os.ModePerm); err != nil {
				return err
			}

			cmd := exec.Command("protector4j", "-t", "java", "-f", "task.java.yaml")

			cmd.Dir = temp

			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout

			if err := cmd.Run(); err != nil {
				return err
			}

			for _, jarPath := range task.JarsPath {
				jar, err := os.Stat(jarPath)
				if err != nil {
					return err
				}

				content, err := ioutil.ReadFile(filepath.Join(temp, jar.Name()))
				if err != nil {
					return err
				}

				if err := ioutil.WriteFile(jarPath, content, jar.Mode()); err != nil {
					return err
				}
			}

			if err := os.RemoveAll(filepath.Join(output, "jre")); err != nil {
				return err
			}
			if err := os.Rename(filepath.Join(temp, "jre"), filepath.Join(output, "jre")); err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
