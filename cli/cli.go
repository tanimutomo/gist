package cli

import (
	"errors"
	"fmt"

	"github.com/b4b4r07/gist/api"
)

func NewGist() (*api.Gist, error) {
	return api.NewGist(api.Config{
		Token:      Conf.Gist.Token,
		BaseURL:    Conf.Gist.BaseURL,
		NewPrivate: Conf.Flag.NewPrivate,
		ClonePath:  Conf.Gist.Dir,
	})
}

func Edit(g *api.Gist, fname string) error {
	if err := g.Sync(fname); err != nil {
		return err
	}

	editor := Conf.Core.Editor
	if editor == "" {
		return errors.New("$EDITOR: not set")
	}

	if err := Run(editor, fname); err != nil {
		return err
	}

	return g.Sync(fname)
}

func Sync(g *api.Gist, fname string) error {
	kind, content, err := g.Compare(fname)
	if err != nil {
		return err
	}
	// TODO
	_ = content
	switch kind {
	case "local":
		err = g.UpdateRemote(fname, content)
		fmt.Printf("Uploaded\t%s\n", fname)
	case "remote":
		err = g.UpdateLocal(fname, content)
		fmt.Printf("Downloaded\t%s\n", fname)
	case "equal":
		fmt.Printf("Not changed\t%s\n", fname)
	case "":
		// Locally but not remote
	default:
	}

	return err
}
