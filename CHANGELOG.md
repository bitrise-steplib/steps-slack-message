## Changelog (Current version: 2.6.3)

-----------------

### 2.6.3 (2017 Oct 30)


### 2.6.2 (2017 Aug 28)

* [9bc5708] Added release_config.yml
* [69965f4] Prepare for 2.6.2
* [a13edaa] check if body sring is ok (#23)

### 2.6.1 (2017 Aug 07)

* [14a9b8c] STEP_VERSION: 2.6.1
* [75c5e8b] Add option for Linkifying user name and channel name (#22)

### 2.6.0 (2017 Aug 07)

* [919a33c] STEP_VERSION: 2.6.0
* [2d687d1] username example wasn't clear enough, actually it was a bit misleading (#20)

### 2.5.0 (2017 May 22)

* [b2615bb] STEP_VERSION: 2.5.0
* [053d908] Make sure that the backslash+n char sequences are converted to the `\n` newline escape char in the message texts (#17)

### 2.4.1 (2017 May 09)

* [9add234] v2.4.1 - ImageURL input title quickfix
* [43bc9cb] bitrise.yml normalized by workflow editor (#16)
* [50b7707] step.yml quickfix - title of image_url_on_error

### 2.4.0 (2017 May 09)

* [c68cef2] Merge pull request #15 from bitrise-io/feature/preps-for-v2.4.0
* [7ee1395] Merge branch 'master' into feature/preps-for-v2.4.0
* [676e074] smaller test gif, which shows in Slack immediately
* [fb4d305] minimal step.yml revision: type tag and host os tags update
* [4c1d1d1] go deps update
* [d7a5d7a] step version bump to 2.4.0
* [48c3530] removed old, unused depman dep files
* [1dfb173] Merge pull request #14 from simonmartyr/master
* [3949228] Test
* [d3966b3] House keeping
* [a5f0bbf] fix missing :
* [4c243aa] yaml fix
* [e97aecc] image within attachment.

### 2.3.0 (2016 Sep 12)

* [e8b0c06] prep for v2.3.0
* [aaeaf50] Merge pull request #10 from bitrise-io/feature/markdown-in-attachment-mode
* [526df8b] added "fields" to markdown list
* [d21a837] markdown formatting support in attachment mode

### 2.2.0 (2016 Sep 12)

* [129d26c] step.yml : specifying bin_name
* [baa97b9] Merge pull request #9 from bitrise-io/feature/slack-msg-color
* [53969d2] prep for v2.2.0
* [1f95495] step.yml : use the new deps syntax
* [6570203] minor check syntax change
* [6e291be] new formatting mode, attachment
* [03dffcf] Merge pull request #8 from bitrise-io/feature/go-toolkit-and-revision
* [a7c51c5] fix in step.sh
* [2c7f206] base revision DONE
* [df3c682] fail test workflow fix
* [7ad39b7] removed explicitely failing test from test
* [9120630] readme update
* [bcef3fb] base revision / update
* [5bec584] Merge pull request #3 from vasarhelyia/patch-1
* [c0c0fab] Update README.md
* [50eccc1] typo fix

### 2.1.0 (2015 Aug 17)

* [756f39f] step.yml typo fix
* [30c8262] revision: inverted IconURL<->emoji priority : if URL defined the Emoji is ignored

### 2.0.0 (2015 Aug 17)

* [74902fe] migrated to the new, V2 syntax

### 1.2.0 (2015 Apr 03)

* [73a01ff] minor description changes
* [a9bb98d] is_always_run is now defaults to 'true'; STEPLIB_BUILD_STATUS based build status (ok/error) handling added - related optional message, username and icon params can be defined separately for the error state
* [88c082a] replaced the old _utils and _formatted_output with the 'depman' based bash-utils dependency

### 1.1.0 (2015 Feb 02)

* [d812a10] core rewrite in Go; added Icon Emoji and URL handling; Channel and Username (From Name) is now not required (as it is not required by Slack)
* [0184c73] utils fix - not relevant (not used)

### 1.0.2 (2015 Jan 16)

* [40d6eec] step.YML update & proper slack message encoding with inline Ruby JSON dump
* [1ceca5b] Merge pull request #1 from erosdome/master
* [1790d96] Update README.md

### 1.0.1 (2014 Oct 30)

* [eccfb7a] removed the 'ghost' icon - now using the default Slack icon for the msg

-----------------

Updated: 2017 Oct 30