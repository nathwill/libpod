package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/projectatomic/libpod/libpod"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	startFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "attach, a",
			Usage: "Attach container's STDOUT and STDERR",
		},
		cli.StringFlag{
			Name:  "detach-keys",
			Usage: "Override the key sequence for detaching a container. Format is a single character [a-Z] or ctrl-<value> where <value> is one of: a-z, @, ^, [, , or _.",
		},
		cli.BoolFlag{
			Name:  "interactive, i",
			Usage: "Keep STDIN open even if not attached",
		},
		cli.BoolFlag{
			Name:  "sig-proxy",
			Usage: "proxy received signals to the process",
		},
		LatestFlag,
	}
	startDescription = `
   podman start

   Starts one or more containers.  The container name or ID can be used.
`

	startCommand = cli.Command{
		Name:                   "start",
		Usage:                  "Start one or more containers",
		Description:            startDescription,
		Flags:                  startFlags,
		Action:                 startCmd,
		ArgsUsage:              "CONTAINER-NAME [CONTAINER-NAME ...]",
		UseShortOptionHandling: true,
	}
)

func startCmd(c *cli.Context) error {
	args := c.Args()
	if len(args) < 1 && !c.Bool("latest") {
		return errors.Errorf("you must provide at least one container name or id")
	}

	attach := c.Bool("attach")

	if len(args) > 1 && attach {
		return errors.Errorf("you cannot start and attach multiple containers at once")
	}

	if err := validateFlags(c, startFlags); err != nil {
		return err
	}

	if c.Bool("sig-proxy") && !attach {
		return errors.Wrapf(libpod.ErrInvalidArg, "you cannot use sig-proxy without --attach")
	}

	runtime, err := getRuntime(c)
	if err != nil {
		return errors.Wrapf(err, "error creating libpod runtime")
	}
	defer runtime.Shutdown(false)
	if c.Bool("latest") {
		lastCtr, err := runtime.GetLatestContainer()
		if err != nil {
			return errors.Wrapf(err, "unable to get latest container")
		}
		args = append(args, lastCtr.ID())
	}
	var lastError error
	for _, container := range args {
		ctr, err := runtime.LookupContainer(container)
		if err != nil {
			if lastError != nil {
				fmt.Fprintln(os.Stderr, lastError)
			}
			lastError = errors.Wrapf(err, "unable to find container %s", container)
			continue
		}

		ctrState, err := ctr.State()
		if err != nil {
			return errors.Wrapf(err, "unable to get container state")
		}

		if attach {
			inputStream := os.Stdin
			if !c.Bool("interactive") {
				inputStream = nil
			}
			if ctrState == libpod.ContainerStateRunning {
				return attachCtr(ctr, os.Stdout, os.Stderr, inputStream, c.String("detach-keys"), c.BoolT("sig-proxy"))
			}

			if err := startAttachCtr(ctr, os.Stdout, os.Stderr, inputStream, c.String("detach-keys"), c.Bool("sig-proxy")); err != nil {
				return errors.Wrapf(err, "unable to start container %s", ctr.ID())
			}

			if ecode, err := ctr.ExitCode(); err != nil {
				logrus.Errorf("unable to get exit code of container %s: %q", ctr.ID(), err)
			} else {
				exitCode = int(ecode)
			}

			return ctr.Cleanup()
		}
		if ctrState == libpod.ContainerStateRunning {
			fmt.Println(ctr.ID())
			continue
		}
		// Handle non-attach start
		if err := ctr.Start(); err != nil {
			if lastError != nil {
				fmt.Fprintln(os.Stderr, lastError)
			}
			lastError = errors.Wrapf(err, "unable to start container %q", container)
			continue
		}
		fmt.Println(ctr.ID())
	}

	return lastError
}
