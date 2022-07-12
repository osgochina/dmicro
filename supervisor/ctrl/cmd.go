package ctrl

import (
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"gopkg.in/AlecAivazis/survey.v1"
	"strings"
	"time"
)

var App = grumble.New(&grumble.Config{
	Name:                  "DServer",
	Description:           "好用的服务管理工具",
	HistoryFile:           "/tmp/foo.hist",
	Prompt:                "DSvr » ",
	PromptColor:           color.New(color.FgGreen, color.Bold),
	HelpHeadlineColor:     color.New(color.FgGreen),
	HelpHeadlineUnderline: true,
	HelpSubCommands:       true,

	Flags: func(f *grumble.Flags) {
		f.String("d", "directory", "DEFAULT", "set an alternative root directory path")
		f.Bool("v", "verbose", false, "enable verbose mode")
	},
})

func init() {
	App.SetPrintASCIILogo(func(a *grumble.App) {
		fmt.Println("  ____    ____                                      ")
		fmt.Println(" |  _ \\  / ___|    ___   _ __  __   __   ___   _ __ ")
		fmt.Println(" | | | | \\___ \\   / _ \\ | '__| \\ \\ / /  / _ \\ | '__|")
		fmt.Println(" | |_| |  ___) | |  __/ | |     \\ V /  |  __/ | |   ")
		fmt.Println(" |____/  |____/   \\___| |_|      \\_/    \\___| |_|  ")
		fmt.Println()
	})
	App.AddCommand(&grumble.Command{
		Name: "flags",
		Help: "test flags",
		Flags: func(f *grumble.Flags) {
			f.Duration("d", "duration", time.Second, "duration test")
			f.Int("i", "int", 1, "test int")
			f.Int64("l", "int64", 2, "test int64")
			f.Uint("u", "uint", 3, "test uint")
			f.Uint64("j", "uint64", 4, "test uint64")
			f.Float64("f", "float", 5.55, "test float64")
		},
		Run: func(c *grumble.Context) error {
			fmt.Println("duration ", c.Flags.Duration("duration"))
			fmt.Println("int      ", c.Flags.Int("int"))
			fmt.Println("int64    ", c.Flags.Int64("int64"))
			fmt.Println("uint     ", c.Flags.Uint("uint"))
			fmt.Println("uint64   ", c.Flags.Uint64("uint64"))
			fmt.Println("float    ", c.Flags.Float64("float"))
			return nil
		},
	})
	promptCommand := &grumble.Command{
		Name: "prompt",
		Help: "set a custom prompt",
	}
	App.AddCommand(promptCommand)

	promptCommand.AddCommand(&grumble.Command{
		Name: "set",
		Help: "set a custom prompt",
		Run: func(c *grumble.Context) error {
			c.App.SetPrompt("CUSTOM PROMPT >> ")
			return nil
		},
	})

	promptCommand.AddCommand(&grumble.Command{
		Name: "reset",
		Help: "reset to default prompt",
		Run: func(c *grumble.Context) error {
			c.App.SetDefaultPrompt()
			return nil
		},
	})

	App.AddCommand(&grumble.Command{
		Name:    "daemon",
		Help:    "run the daemon",
		Aliases: []string{"run"},
		Flags: func(f *grumble.Flags) {
			f.Duration("t", "timeout", time.Second, "timeout duration")
		},
		Args: func(a *grumble.Args) {
			a.Bool("production", "whether to start the daemon in production or development mode")
			a.Int("opt-level", "the optimization mode", grumble.Default(3))
			a.StringList("services", "additional services that should be started", grumble.Default([]string{"test", "te11"}))
		},
		Run: func(c *grumble.Context) error {
			c.App.Println("timeout:", c.Flags.Duration("timeout"))
			c.App.Println("directory:", c.Flags.String("directory"))
			c.App.Println("verbose:", c.Flags.Bool("verbose"))
			c.App.Println("production:", c.Args.Bool("production"))
			c.App.Println("opt-level:", c.Args.Int("opt-level"))
			c.App.Println("services:", strings.Join(c.Args.StringList("services"), ","))
			return nil
		},
	})

	App.AddCommand(&grumble.Command{
		Name: "args",
		Help: "test args",
		Args: func(a *grumble.Args) {
			a.String("s", "test string")
			a.Duration("d", "test duration", grumble.Default(time.Second))
			a.Int("i", "test int", grumble.Default(5))
			a.Int64("i64", "test int64", grumble.Default(int64(-88)))
			a.Uint("u", "test uint", grumble.Default(uint(66)))
			a.Uint64("u64", "test uint64", grumble.Default(uint64(8888)))
			a.Float64("f64", "test float64", grumble.Default(float64(5.889)))
			a.StringList("sl", "test string list", grumble.Default([]string{"first", "second", "third"}), grumble.Max(3))
		},
		Run: func(c *grumble.Context) error {
			fmt.Println("s  ", c.Args.String("s"))
			fmt.Println("d  ", c.Args.Duration("d"))
			fmt.Println("i  ", c.Args.Int("i"))
			fmt.Println("i64", c.Args.Int64("i64"))
			fmt.Println("u  ", c.Args.Uint("u"))
			fmt.Println("u64", c.Args.Uint64("u64"))
			fmt.Println("f64", c.Args.Float64("f64"))
			fmt.Println("sl ", strings.Join(c.Args.StringList("sl"), ","))
			return nil
		},
	})
	adminCommand := &grumble.Command{
		Name:     "admin",
		Help:     "admin tools",
		LongHelp: "super administration tools",
	}
	App.AddCommand(adminCommand)

	adminCommand.AddCommand(&grumble.Command{
		Name: "root",
		Help: "root the machine",
		Run: func(c *grumble.Context) error {
			fmt.Println(c.Flags.String("directory"))
			return fmt.Errorf("failed")
		},
	})

	adminCommand.AddCommand(&grumble.Command{
		Name: "kill",
		Help: "kill the process",
		Run: func(c *grumble.Context) error {
			return fmt.Errorf("failed")
		},
	})
	App.AddCommand(&grumble.Command{
		Name: "ask",
		Help: "ask the user for foo",
		Run: func(c *grumble.Context) error {
			ask()
			return nil
		},
	})
}

// the questions to ask
var qs = []*survey.Question{
	{
		Name:      "name",
		Prompt:    &survey.Input{Message: "What is your name?"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name: "color",
		Prompt: &survey.Select{
			Message: "Choose a color:",
			Options: []string{"red", "blue", "green"},
			Default: "red",
		},
	},
	{
		Name:   "age",
		Prompt: &survey.Input{Message: "How old are you?"},
	},
}

func ask() {
	password := ""
	prompt := &survey.Password{
		Message: "Please type your password",
	}
	survey.AskOne(prompt, &password, nil)

	// the answers will be written to this struct
	answers := struct {
		Name          string // survey will match the question and field names
		FavoriteColor string `survey:"color"` // or you can tag fields to match a specific name
		Age           int    // if the types don't match exactly, survey will try to convert for you
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%s chose %s.", answers.Name, answers.FavoriteColor)
}
