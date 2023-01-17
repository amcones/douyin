# 代码贡献流程
## 贡献流程
本项目采用 Github Forking 工作流。

Github Forking 工作流详细步骤如下：

1. Fork 远程仓库到自己的账号下

访问 https://github.com/amcones/douyin ， 点击右上角 [**Fork** 按钮](https://github.com/amcones/douyin/fork)。 fork 后的仓库地址为：https://github.com/your-name/douyin

2. 克隆 fork 的仓库到本地

```shell
$ git clone https://github.com/your-name/douyin
$ cd douyin
$ git remote add upstream https://github.com/amcones/douyin
$ git remote set-url --push upstream no_push # Never push to upstream develop
$ git remote -v # Confirm that your remotes make sense
origin https://github.com/your-name/douyin (fetch)
origin https://github.com/your-name/douyin (push)
upstream https://github.com/amcones/douyin (fetch)
upstream https://github.com/amcones/douyin (push)
```

3. 创建功能分支

首先同步本地仓库的 main 分支为最新的状态（跟 upstream main 分支一致）

```shell
$ git fetch upstream
$ git checkout main
$ git rebase upstream/main
```

创建功能分支：

```shell
$ git checkout -b feature/add-function
```

4. 开发代码

先 rebase 到 `upstream/main`:

```shell
$ git fetch upstream # commit 前需要再次同步 feature 跟 upstream/main
$ git rebase upstream/main
```

5. 提交 commit

在 feature/add-function 分支上开发代码，开发完代码后，提交commit：

```shell
$ git add .
$ git status
$ git commit
```

分支开发完成后，可能会有一堆 commit，但是合并到主干的时候，往往希望只有一个（或最多两三个）commit，这样不仅清晰，也容易管理。可以用 git rebase 来合并修改我们的 commit。

```shell
$ git rebase -i origin/main
```

`-i` 参数表示进入交互模式，例如：

```shell
pick 07c5abd Introduce OpenPGP and teach basic usage
pick de9b1eb Fix PostChecker::Post#urls
pick 3e7ee36 Hey kids, stop all the highlighting
pick fa20af3 git interactive rebase, squash, amend

# Rebase 8db7e8b..fa20af3 onto 8db7e8b
#
# Commands:
#  p, pick = use commit
#  r, reword = use commit, but edit the commit message
#  e, edit = use commit, but stop for amending
#  s, squash = use commit, but meld into previous commit
#  f, fixup = like "squash", but discard this commit's log message
#  x, exec = run command (the rest of the line) using shell
#
# These lines can be re-ordered; they are executed from top to bottom.
#
# If you remove a line here THAT COMMIT WILL BE LOST.
#
# However, if you remove everything, the rebase will be aborted.
#
# Note that empty commits are commented out
```

上面的交互界面中，先列出当前分支最新的 4 个 commit（越下面越新）。每个 commit 前面有一个操作命令，默认是 pick，表示该行 commit 被选中，要进行 rebase 操作。
4 个 commit 的下面是一大堆注释，列出可以使用的命令：
- pick：正常选中
- reword：选中，并且修改提交信息；
- edit：选中，rebase 时会暂停，允许你修改这个 commit（参考这里）
- squash：选中，会将当前 commit 与上一个 commit 合并
- fixup：与 squash 相同，但不会保存当前 commit 的提交信息
- exec：执行其他 shell 命令

上面这 6 个命令当中，squash 和 fixup 可以用来合并 commit。先把需要合并的 commit 前面的动词，改成 squash（或者 s）：

```shell
pick 07c5abd Introduce OpenPGP and teach basic usage
s de9b1eb Fix PostChecker::Post#urls
s 3e7ee36 Hey kids, stop all the highlighting
pick fa20af3 git interactive rebase, squash, amend
```

这样一改，执行后，当前分支只会剩下两个 commit。第二行和第三行的 commit，都会合并到第一行的 commit。提交信息会同时包含，这三个 commit 的提交信息：

```shell
# This is a combination of 3 commits.
# The first commit's message is:
Introduce OpenPGP and teach basic usage

# This is the 2nd commit message:
Fix PostChecker::Post#urls

# This is the 3rd commit message:
Hey kids, stop all the highlighting
```

如果将第三行的 squash 命令改成 fixup 命令：

```shell
pick 07c5abd Introduce OpenPGP and teach basic usage
s de9b1eb Fix PostChecker::Post#urls
f 3e7ee36 Hey kids, stop all the highlighting
pick fa20af3 git interactive rebase, squash, amend
```

运行结果相同，还是会生成两个 commit，第二行和第三行的 commit，都合并到第一行的 commit。但是，新的提交信息里面，第三行 commit 的提交信息，会被注释掉。

```shell
# This is a combination of 3 commits.
# The first commit's message is:
Introduce OpenPGP and teach basic usage

# This is the 2nd commit message:
Fix PostChecker::Post#urls

# This is the 3rd commit message:
# Hey kids, stop all the highlighting
```

还有另外一种合并 commit 的简便方法，就是先撤销过去 5 个 commit，然后再建一个新的：

```shell
$ git reset HEAD~5
$ git add .
$ git commit -am "Here's the bug fix that closes #28"
$ git push --force
```

squash 和 fixup 命令，还可以当作命令行参数使用，自动合并 commit。

```shell
$ git commit --fixup
$ git rebase -i --autosquash
```

6) push 功能分支到个人远程仓库

在完成了开发，并 commit 后，需要将功能分支 push 到远程代码仓库：

```shell
$ git push -f origin feature/add-function
```

7) 在个人远程仓库页面创建 pull request。

提交到远程仓库以后，就可以发出 pull request 到 main 分支，后面由 团队的 Maintainer 进行代码 Review，并 merge 到主干代码。

## Commit规范

本项目使用知名的[Angular commit message规范](https://github.com/angular/angular/blob/main/CONTRIBUTING.md#-commit-message-format)。

其中比较重要的是**Type**:

- build: Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)
- ci: Changes to our CI configuration files and scripts (examples: CircleCi, SauceLabs)
- docs: Documentation only changes
- feat: A new feature
- fix: A bug fix
- perf: A code change that improves performance
- refactor: A code change that neither fixes a bug nor adds a feature
- test: Adding missing tests or correcting existing tests

在**scope**上，不做强行的规定，一般有:

- GLOBAL
- controller
- service
- entity
- config
- api
- doc
