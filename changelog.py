import git
import os
import sys
import re

def get_remote_url():
    """Fetch the GitHub remote repository URL"""
    try:
        repo = git.Repo(os.getcwd(), search_parent_directories=True)
        remote_url = repo.remotes.origin.url  # Get remote URL

        # Convert SSH URL (git@github.com:user/repo.git) to HTTPS
        if remote_url.startswith("git@github.com:"):
            remote_url = re.sub(r"git@github\.com:(.*)\.git", r"https://github.com/\1", remote_url)
        elif remote_url.endswith(".git"):
            remote_url = remote_url[:-4]  # Remove .git if it's there

        return remote_url
    except Exception:
        return None  # Return None if no remote is found

def get_all_commits(remote_url):
    """Fetch all commit details from a local Git repository"""
    try:
        repo = git.Repo(os.getcwd(), search_parent_directories=True)
        commits = list(repo.iter_commits())  # Get all commits
        commit_list = []

        for commit in commits:
            commit_info = {
                "sha": commit.hexsha[:7],  # Short hash
                "author": commit.author.name,
                "date": commit.committed_datetime.strftime("%Y-%m-%d"),
                "message": commit.message.strip(),
                "url": f"{remote_url}/commit/{commit.hexsha}" if remote_url else None,
            }
            commit_list.append(commit_info)

        return commit_list

    except git.exc.InvalidGitRepositoryError:
        print("Error: Not inside a Git repository. Please run this script inside a Git project.")
        sys.exit(1)
    except Exception as e:
        print(f"Unexpected Error: {e}")
        sys.exit(1)

def generate_changelog(remote_url, commits):
    """Format commit history into a changelog"""
    changelog = "# Changelog\n\n"
    if remote_url:
        changelog += f"🔗 **GitHub Repository:** [{remote_url}]({remote_url})\n\n"

    for commit in commits:
        commit_link = f"[View Commit]({commit['url']})" if commit["url"] else ""
        changelog += f"""
## [{commit['sha']}] - {commit['date']}
**Author:** {commit['author']}

### Changes:
- {commit['message']}

{commit_link}

"""
    return changelog.strip()

if __name__ == "__main__":
    remote_url = get_remote_url()  # Get remote URL
    commits = get_all_commits(remote_url)  # Get all commits
    changelog = generate_changelog(remote_url, commits)  # Format changelog
    print(changelog)  # Print to terminal

    # Save to file
    with open("CHANGELOG.md", "w") as file:
        file.write(changelog)

    print("\n✅ Changelog saved to 'CHANGELOG.md'")
