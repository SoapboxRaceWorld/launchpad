package launchpad

import (
	"fmt"
	"github.com/google/go-github/github"
	"strings"
	"github.com/mcuadros/go-version"
)

func CheckForUpdates(requestVersion string) LauncherUpdatePacket {
	packet := LauncherUpdatePacket{}

	releases, _, err := ghClient.Repositories.ListReleases(
		clientContext,
		"SoapboxRaceWorld",
		"GameLauncher_NFSW", &github.ListOptions{PerPage: 100})

	if err != nil {
		fmt.Println(err)
		packet.Code = 1001

		return packet
	}

	latestRelease := releases[0]

	payload := LauncherUpdatePayload{}
	payload.ClientVersion = requestVersion
	payload.LatestVersion = *latestRelease.TagName

	if version.Compare(requestVersion, *latestRelease.TagName, "<") {
		payload.UpdateExists = true
		update := LauncherUpdate{}
		update.DownloadURL = *latestRelease.Assets[0].BrowserDownloadURL
		payload.Update = update
	}

	downloadCounts := LauncherDownloadCounts{}

	downloadCounts.CurrentVersion = *latestRelease.Assets[0].DownloadCount

	for _, release := range releases {
		if *release.Prerelease {
			continue
		}
		downloadCounts.Total += *release.Assets[0].DownloadCount
	}

	payload.DownloadCounts = downloadCounts

	packet.Payload = payload

	return packet
}

func GetLatestChangelog() string {
	releases, _, err := ghClient.Repositories.ListReleases(
		clientContext,
		"SoapboxRaceWorld",
		"GameLauncher_NFSW", &github.ListOptions{PerPage: 100})

	if err != nil {
		fmt.Println(err)
		return "Error occurred, please try again later"
	}

	latestRelease := releases[0]

	return strings.Replace(*latestRelease.Body, "##### CHANGELOG:\r\n", "", 1)
}