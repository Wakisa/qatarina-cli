# QATARINA CLI Developer Docs

## Authentication
Before using commands, you must log in:
    $ qatarina-cli login --email user@example.com --password secret123
Replace user@example.com and secret123 with our actual email and password
To log out:
    $ qatarina-cli logout

# Projects Commands

## Create a Project
We have two ways to do this, you can use flags or interactive wizard
    $qatarina-cli project create \
  --name "My Project" \
  --description "Web app testing" \
  --version "1.0.0" \
  --website-url "https://example.com" \
  --github-url "https://github.com/example/repo"
                OR
If flags are missing, an interactive wizard will launch

## List Projects
    $ qatarina-cli project list

## View Project
    $ qatarina-cli project view 1

## Delete Project
    $ qatarina-cli project delete 1

# Test Cases Commands

## Create Test Case
We have two ways here, using flags or interactive wizard
    $ qatarina-cli test-case create \
  --title "Login with valid credentials" \
  --kind "general" \
  --project 1 \
  --description "Ensure login works with correct credentials" \
  --code "TC-001" \
  --feature-or-module "Authentication" \
  --tags "login,smoke" \
  --draft=false
            OR 
If required flags are missing, an interactive wizard will launch.

## List Test Cases
    $ qatarina-cli test-case list --project 1

## Update Test Case
     $ qatarina-cli test-case update 10 \
  --title "Login with valid credentials (updated)" \
  --draft=true

## Delete Test Cases
    $ qatarina-cli test-case delete 10

## Import Test Cases (Excel/CSV)
You can bulk import test cases from an Excel (.xlsx)
    $ qatarina-cli import-file --project 1 --file ./testcases.xlsx
File must contain headers: check files in testdata/
    e.g Title | Description | Kind | Code | FeatureOrModule | Tags | IsDraft

# Test Plan Commnads

## Assign Test Cases to a Test Plan
You can assign existing test cases from a project into a specific test plan.
    $ qatarina-cli assign-cases --project <projectID> --plan <planID>
This command launches an interactive wizard where you can select which test cases to assign. Doing this you will be also assigning a testcase to a user(use commas if many).
Use SPACEBAR to select and ENTER to confirm.

# Modules Commands

## List Modules in a Project
    $ qatarina-cli project modules 1

## Create Module
    $ qatarina-cli module create \
  --project-id 1 \
  --name "Authentication" \
  --code "AUTH" \
  --priority 1 \
  --type "feature" \
  --description "Handles user login/logout"

## Update Module
    $ qatarina-cli module update 5 \
  --name "Auth System" \
  --priority 2

## List Modules
    $ qatarina-cli module list

## Delete Module
    $ qatarina-cli module delete 5