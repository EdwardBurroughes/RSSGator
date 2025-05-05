package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/EdBurroughes/rss-blog-aggregator/internal/config"
	"github.com/EdBurroughes/rss-blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

func buildCommandFromArgs() (command, error) {
	arguments := os.Args
	if len(arguments) < 2 {
		return command{}, fmt.Errorf("uh oh missing expected arguments")
	}
	return command{
		name:      arguments[1],
		arguments: arguments[2:],
	}, nil

}

func main() {
	rssConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", rssConfig.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	if err != nil {
		log.Fatal(err)
	}
	s := state{
		dbQueries,
		&rssConfig,
	}
	initCmds := make(map[string]func(*state, command) error)
	cmds := commands{
		initCmds,
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handleReset)
	cmds.register("users", handleGetUsers)
	cmds.register("agg", handleAgg)
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("feeds", handleFeeds)
	cmds.register("follow", middlewareLoggedIn(handleFollow))
	cmds.register("following", middlewareLoggedIn(handleFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollowing))
	cmds.register("browse", middlewareLoggedIn(handleBrowse))
	cmd, err := buildCommandFromArgs()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmds.run(&s, cmd); err != nil {
		log.Fatal(err)
	}
}
