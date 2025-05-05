package subscriptions

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/thrasher-corp/gocryptotrader/exchanges/subscription"
	"golang.org/x/term"
)

// Public errors
var ErrConflict = errors.New("Subscription configuration conflict")

var (
	errUserInput  = errors.New("invalid user chocie")
	errMaxRetries = errors.New("max input retries exceeded")
)

type Mutator func(subs subscription.List) (subscription.List, error)

func Upgrade(e []byte, oldDefaults subscription.List, newDefaults subscription.List, fn Mutator, desc string) ([]byte, error) {
	newSubsJSON, err := json.Marshal(newDefaults)
	if err != nil {
		return e, err
	}

	s, _, _, err := jsonparser.Get(e, "features", "subscriptions")
	if err != nil {
		if errors.Is(err, jsonparser.KeyPathNotFoundError) {
			return jsonparser.Set(e, newSubsJSON, "features", "subscriptions")
		}
		return e, err
	}

	userSubs := subscription.List{}
	err = json.Unmarshal(s, &userSubs)
	if err != nil {
		return e, err
	}

	if len(userSubs) == 0 {
		return jsonparser.Set(e, newSubsJSON, "features", "subscriptions")
	}

	fmt.Println("Subscription Update Summary:", "---------------------------", desc)

	oldDefaultsStore, err := subscription.NewStoreFromList(oldDefaults)
	if err != nil {
		return e, err
	}

	userAdded, userRemoved := oldDefaultsStore.Diff(userSubs)
	if len(userAdded) == 0 && len(userRemoved) == 0 {
		// User didn't change the original defaults
		return jsonparser.Set(e, newSubsJSON, "features", "subscriptions")
	}

	newDefaultsStore, err := subscription.NewStoreFromList(newDefaults)
	if err != nil {
		return e, err
	}

	added, removed := newDefaultsStore.Diff(userSubs)

	if len(added) == 0 && len(removed) == 0 {
		// User already changed config to match new defaults
		return e, nil
	}

	// TODO: Invoke mutator; DO NOT COMMIT

	eName, err := jsonparser.GetString(e, "name")
	if err != nil {
		return e, err
	}

	fmt.Printf(`Conflict during subscription upgrade for %s.
Default subscriptions have changed, but also editted in config.

Differences in config from old defaults:
---------------------------------------
%s
Differences from new default subscriptions:
------------------------------------------
%s
`, eName, listDiff(userAdded, userRemoved), listDiff(added, removed))

	action, err := promptForAction()
	if err != nil {
		return e, err
	}

	switch action {
	case "A": // Accept the new defaults
		return jsonparser.Set(e, newSubsJSON, "features", "subscriptions")
	case "K": // Keep the existing configuration
		return e, nil
	}

	return e, ErrConflict // Quit and leave the conflict in place
}

func promptForAction() (string, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return "Q", nil
	}

	for range 3 {
		fmt.Println("[Q]uit to edit config", "[A]ccept new Defaults", "[K]eep existing")
		choice := ""
		if _, err := fmt.Scanln(&choice); err != nil {
			return "", fmt.Errorf("%w: %w", errUserInput, err)
		}
		choice = strings.ToUpper(choice)
		switch choice {
		case "Q", "K", "A":
			return choice, nil
		default:
			fmt.Println("Unknown response; Please enter Q, A, or K followed by return")
		}
	}
	return "", fmt.Errorf("%w: %w", errUserInput, errMaxRetries)
}

func listDiff(added, removed subscription.List) (diff string) {
	if len(added) > 0 {
		diff = "+ " + strings.Join(added.Strings(), "\n+ ") + "\n"
	}
	if len(removed) > 0 {
		diff = diff + "- " + strings.Join(removed.Strings(), "\n- ") + "\n"
	}
	return diff
}
