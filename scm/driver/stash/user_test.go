// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jenkins-x/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/h2non/gock.v1"
)

func TestUserFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://bitbucket.example.com").
		Get("plugins/servlet/applinks/whoami").
		Reply(200).
		Type("text/plain").
		BodyString("jcitizen")

	gock.New("https://bitbucket.example.com").
		Get("rest/api/1.0/users/jcitizen").
		Reply(200).
		Type("application/json").
		File("testdata/user.json")

	client, _ := New("https://bitbucket.example.com")
	got, _, err := client.Users.Find(context.Background())
	if err != nil {
		t.Error(err)
	}

	want := new(scm.User)
	raw, _ := ioutil.ReadFile("testdata/user.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestUserLoginFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://bitbucket.example.com").
		Get("rest/api/1.0/users/jcitizen").
		Reply(200).
		Type("application/json").
		File("testdata/user.json")

	client, _ := New("https://bitbucket.example.com")
	got, _, err := client.Users.FindLogin(context.Background(), "jcitizen")
	if err != nil {
		t.Error(err)
	}

	want := new(scm.User)
	raw, _ := ioutil.ReadFile("testdata/user.json.golden")
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestUserFindEmail(t *testing.T) {
	defer gock.Off()

	gock.New("https://bitbucket.example.com").
		Get("plugins/servlet/applinks/whoami").
		Reply(200).
		Type("text/plain").
		BodyString("jcitizen")

	gock.New("https://bitbucket.example.com").
		Get("rest/api/1.0/users/jcitizen").
		Reply(200).
		Type("application/json").
		File("testdata/user.json")

	client, _ := New("https://bitbucket.example.com")
	email, _, err := client.Users.FindEmail(context.Background())
	if err != nil {
		t.Error(err)
	}

	if got, want := email, "jane@example.com"; got != want {
		t.Errorf("Want email %s, got %s", want, got)
	}
}

func TestUserListInvitations(t *testing.T) {
	client, _ := New("https://bitbucket.example.com")

	invites, res, err := client.Users.ListInvitations(context.Background())

	assert.Len(t, invites, 0, "Should be an empty list")
	assert.Equal(t, res.Status, 200, "Should be success response")
	assert.NoError(t, err, "Should not be an error")
}

func TestUserAcceptInvite(t *testing.T) {
	client, _ := New("https://bitbucket.example.com")

	res, err := client.Users.AcceptInvitation(context.Background(), 0)

	assert.Equal(t, res.Status, 200, "Should be success response")
	assert.NoError(t, err, "Should not be an error")
}
