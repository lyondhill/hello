package main

import (
	"fmt"
	"github.com/shazow/go-git"
)

func main() {
	repo, err := git.OpenRepository(".")
	fmt.Println(repo, err)
	commit, err := repo.GetCommitOfBranch("master")
	fmt.Println(commit, err)
	tree, err := repo.GetTree(commit.TreeId().String())
	fmt.Println(tree, err)
	scanner, err := tree.Scanner()
	for scanner.Scan() {
		fmt.Println(scanner.TreeEntry())
	}
}