// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package git

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// CommitFromReader will generate a Commit from a provided reader
// We need this to interpret commits from cat-file or cat-file --batch
//
// If used as part of a cat-file --batch stream you need to limit the reader to the correct size
func CommitFromReader(gitRepo *Repository, objectID ObjectID, reader io.Reader) (*Commit, error) {
	commit := &Commit{
		ID:        objectID,
		Author:    &Signature{},
		Committer: &Signature{},
	}

	payloadSB := new(strings.Builder)
	signatureSB := new(strings.Builder)
	messageSB := new(strings.Builder)
	message := false
	pgpsig := false

	bufReader, ok := reader.(*bufio.Reader)
	if !ok {
		bufReader = bufio.NewReader(reader)
	}

	// Git parses tree (1), parent (0-N), author (1), and committer (1) in that order.
	// Headers outside that position must not populate those fields. A crafted commit
	// can otherwise make Forgejo show a different tree or diff than git does.
	// https://github.com/git/git/blob/fb255012acbcbe81189cb8014da0280b09c66ecc/fsck.c#L969-L999
	type parseState int
	const (
		treeNext parseState = iota
		parentOrAuthorNext
		committerNext
		otherNext
	)
	state := treeNext

readLoop:
	for {
		line, err := bufReader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				if message {
					_, _ = messageSB.Write(line)
				}
				_, _ = payloadSB.Write(line)
				break readLoop
			}
			return nil, err
		}
		if pgpsig {
			if len(line) > 0 && line[0] == ' ' {
				_, _ = signatureSB.Write(line[1:])
				continue
			}
			pgpsig = false
		}

		if !message {
			// This is probably not correct but is copied from go-gits interpretation...
			trimmed := bytes.TrimSpace(line)
			if len(trimmed) == 0 {
				message = true
				_, _ = payloadSB.Write(line)
				continue
			}

			split := bytes.SplitN(trimmed, []byte{' '}, 2)
			var data []byte
			if len(split) > 1 {
				data = split[1]
			}

			noun := string(split[0])
			switch state {
			case treeNext:
				if noun == "tree" {
					commit.Tree = *NewTree(gitRepo, MustIDFromString(string(data)))
					state = parentOrAuthorNext
				} else {
					return nil, fmt.Errorf("CommitFromReader: expected tree, but received %q", noun)
				}
				_, _ = payloadSB.Write(line)
			case parentOrAuthorNext:
				switch noun {
				case "parent":
					commit.Parents = append(commit.Parents, MustIDFromString(string(data)))
				case "author":
					commit.Author = &Signature{}
					commit.Author.Decode(data)
					state = committerNext
				default:
					return nil, fmt.Errorf("CommitFromReader: expected parent or author, but received %q", noun)
				}
				_, _ = payloadSB.Write(line)
			case committerNext:
				if noun == "committer" {
					commit.Committer = &Signature{}
					commit.Committer.Decode(data)
					state = otherNext
				} else {
					return nil, fmt.Errorf("CommitFromReader: expected committer, but received %q", noun)
				}
				_, _ = payloadSB.Write(line)
			case otherNext:
				switch noun {
				case "gpgsig":
					fallthrough
				case "gpgsig-sha256": // FIXME: no intertop, so only 1 exists at present.
					_, _ = signatureSB.Write(data)
					_ = signatureSB.WriteByte('\n')
					pgpsig = true
				default:
					// encoding, change-id, and any later tree/parent/author/committer
					// stay in the payload and do not replace the positional headers.
					_, _ = payloadSB.Write(line)
				}
			}
		} else {
			_, _ = messageSB.Write(line)
			_, _ = payloadSB.Write(line)
		}
	}
	commit.CommitMessage = messageSB.String()
	commit.Signature = &ObjectSignature{
		Signature: signatureSB.String(),
		Payload:   payloadSB.String(),
	}
	if len(commit.Signature.Signature) == 0 {
		commit.Signature = nil
	}

	return commit, nil
}
