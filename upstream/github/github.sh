#!/bin/bash
#
# This github driver do a lot of things:
#
# For the infra repo (identified by --infra option) do the following
#
# - ensure the infra is a git repo.
# - ensure it is attached to a remote git repo. It will create anything if needed.
# - update infra through git commit and push to master.
#
# TODO: re-write it in GO.
# This code is very basic, and considered as POC for forjj
# It should be re-wrote as a go program instead 
# to manipulate yaml files, API call and strings easily.

SOUT='{ "repos": [%s],
  "services": [%s],
  "state_code": %s,
  "status": "%s" }'

FOUT='{ "state_code": %s, "error_message": "%s" }'
SOUT_REPO='{
     "name": "%s",
     "upstream": "%s",
     "config": "%s/%s"
     }'
SOUT_SERVICE='{
      "upstream": "%s"
     }'


declare -A GITHUB


GITHUB[DEBUG]=0

function parse_options
{ # Parse and identify forjj core options.
 while [ $# -ne 0 ]
 do
   case "$1" in
     --forjj-infra)
       INFRA=$2; shift;shift;;
     --forjj-organization)
       ORGA=$2; shift;shift;;
     --infra)
       INFRA_REPOS=true; shift;;
     *)
       if [ "${1:0:9}" != "--github-" ]
       then
         shift
         continue
       fi
       TYPE=${1:9}
       GITHUB[${TYPE^^[a-z]}]="$2"
       shift;shift;;
   esac
 done
}

function debug_off
{
 if [[ $- =~ x ]]
 then
    set +x
    DEBUG_SAVED=true
 else
    DEBUG_SAVED=false
 fi
}

function debug_on
{
 if [ $DEBUG_SAVED = true ]
 then
    DEBUG_SAVED=false
    set -x
 fi
}

function github-get
{
 debug_off
 local RES
 local CODE
 if [ "${GITHUB[DEBUG]}" -ge 1 ]
 then
    echo "GET: ${GITHUB[api_url]}$1" 1>&2
 fi
 RES="$(curl -s -w 'RETURN_CODE: %{http_code}' -H "$HEADER" ${GITHUB[api_url]}$1)"
 echo "$RES" | grep -v RETURN_CODE:
 CODE=$(echo "$RES" | awk '$1 ~ /RETURN_CODE:/ { print $2 }')
 if [ "${GITHUB[DEBUG]}" -ge 2 ]
 then
    echo "GET RET: $CODE\n$RES" 1>&2
 fi
 debug_on
 return $CODE

}

function github-post
{
 debug_off
 local RES
 if [ "${GITHUB[DEBUG]}" -ge 1 ]
 then
    echo "$POST: ${GITHUB[api_url]}$1\n$2" 1>&2
 fi
 RES="$(curl -s -w "RETURN_CODE: %{http_code}" -H "$HEADER" -X POST -d "$2" ${GITHUB[api_url]}$1)"
 echo "$RES" | grep -v RETURN_CODE:
 CODE=$(echo "$RES" | awk '$1 ~ /RETURN_CODE:/ { print $2 }')
 if [ "${GITHUB[DEBUG]}" -ge 2 ]
 then
    echo "POST RET: $CODE\n$RES" 1>&2
 fi
 debug_on
 return $CODE
}

function json-get
{
 debug_off
 local RES="$1"
 local FILTER="$2"
 shift;shift
 echo "$RES" | jq -r "$(printf "$FILTER" "$@")"
 debug_on
}

function json-test
{
 debug_off
 local RES="$1"
 local FILTER="$2"
 shift;shift
 [ $(echo "$RES" | jq -r "$(printf "$FILTER" "$@")") = true ]
 local RET=$?
 debug_on
 return $RET
}


function create_organization
{
 local RES RET DATA

 if [ "${GITHUB[SERVER]}" = "github.com" ]
 then
    echo "Sorry, but it is not possible to create an organization in the public GITHUB through API. You must create it manually."
    exit
 fi
 RES="$(github-get /organizations)"
 if json-test "$RES" 'map(.login | test("'"$1"'")) | any'
 then
    echo "Organization name '$1' is already used. You need to change the organization name to create a new one or ask to be invited to it."
    exit 1
 fi
 DATA="$(printf '{ "login": "%s", "admin": "%s" }' $1 ${GITHUB[login]})"
 RES="$(github-post /admin/organizations "$DATA")"
 RET=$?
 if [ $RET -ne 201 ]
 then
    echo "There is an issue in the Organization creation."
 fi
}

function ensure_organization_repo
{ # It returns the organization forked repository to use.
 local RES RET DATA title

 # check/create organization
 RES="$(github-get /user/orgs)"
 if json-test "$RES" 'map(.login | test("'"$1"'")) | any'
 then
   echo "No organization '$1' found for user '${GITHUB[login]}'" 1>&2
   if [ "${GITHUB[ORGA-CREATE]}" = true ]
   then
      create_organization "$1"
   else
      echo "To automatically create this organization, add --github-orga-create true. Otherwise, you should create it manually from ${GITHUB[SERVER]}.
If you need to use the existing already created repository, ask the owner to invite you."
   fi
 fi

 # Get Organization Repository
 RES="$(github-get /repos/$1/$2)"
 RET=$?

 if [ $RET -ne 200 ]
 then
    if [ "${GITHUB[${REPO}-TITLE]}" = "" ]
    then
       title="$1 $2 repository automatically created by Forjj(POC)"
    else
       title="${GITHUB[${REPO}-TITLE]}"
    fi
    # Create organization repository.
    DATA="$(printf '{ "name": "%s", "description": "%s", "private": false, 
"has_issues": true, "has_wiki": true, "has_downloads": true }' $REPO "$title")"
    RES="$(github-post /orgs/$1/repos "$DATA")"
    RET=$?
    if [ $RET -ne 201 ]
    then
       printf "The repository has not been created. \nGithub response: %s\n" "$(json-get "$RES" '.message')"
       exit 1
    fi
    GITHUB[new-repo]=true
 else
    GITHUB[new-repo]=false
 fi
 
 # Get user fork repo
 RES=$(github-get /repos/$1/$2/forks)
 RET=$?

 if [ "$(json-get "$RET" 'map(select(.owner.login == "%s"))' "${GITHUB[login]}")" = "" ]
 then
   echo "Forking '$1/$2'..."
   RES="$(github-post /repos/$1/$2/forks '{}')"
   RET=$?
   if [ $? -ne 202 ]
   then
      printf "Unable to fork '%s/%s'. Exiting.\n %s\n" $1 $2 "$(json-get "$RES" '.message')"
      exit 1
   fi
 fi
 GITHUB[repo_name]="$(json-get "$RET" '.name')"
}

function create_service
{
 ensure_organization_repo "$ORGA" "$REPO"
}

#else
#  echo "Warning!! You did not specified any organization to centralize/protect your repos. You are going to use your personal repository as central place." 1>&2
#  GITHUB[repo_name]="$REPO"

  # Ensure the repository exist.
#  RES="$(github-get /repos/${GITHUB[login]}/$REPO)"
#  RET=$?
 
#  if [ $RET -ne 200 ]
#  then
#     if [ "${GITHUB[${REPO}-TITLE]}" = "" ]
#     then
#        title="$ORGA $REPO repository automatically created by Forjj(POC)"
#     else
#        title="${GITHUB[${REPO}-TITLE]}"
#     fi
#     echo "The infra repo '$REPO' doesn't exist. Going to create it."
#     DATA="$(printf '{ "name": "%s", "description": "%s", "private": false, 
#"has_issues": true, "has_wiki": true, "has_downloads": true }' $REPO "$title")"
#     RES="$(github-post /orgs/$ORGA/repos "$DATA")"
#     RET=$?
#     if [ $RET -ne 201 ]
#     then
#        echo "The repository has not been created. $(echo "$RES" | grep message)"
#        exit 1
#     fi
#     echo "Congratulations! Repository '$REPO' has been created in '$FULL_NAME'."
#     GITHUB[new-repo]=true
#  else
#     FULL_NAME="$(json-get "$RES" '.full_name')"
#     echo "Found '$FULL_NAME' repository."
#     GITHUB[new-repo]=false
#  fi
#  GITHUB[repo_name]="$(json-get "$RES" '.name')"
#fi

# if [ $INFRA_REPOS = true ]
# then
#   FORK_SSH_LIKE="ssh://git@${GITHUB[SERVER]}/$LOGIN/$REPO"
#   FORK_SCP_LIKE="git@${GITHUB[SERVER]}:$LOGIN/$REPO"
#   FORK_HTTP_LIKE="https://${GITHUB[SERVER]}/$LOGIN/$REPO"
#
#   ORG_SSH_LIKE="ssh://git@${GITHUB[SERVER]}/$ORGA/$REPO"
#   ORG_SCP_LIKE="git@${GITHUB[SERVER]}:$ORGA/$REPO"
#   ORG_HTTP_LIKE="https://${GITHUB[SERVER]}/$ORGA/$REPO"

#   set -e
#   if [ "${GITHUB[new-repo]}" = true ]
#   then # We need to populate the new repo with any existing data.
#     if [ ! -d /devops/$REPO ]
#     then
#        git init /devops/$REPO
#        cd /devops/$REPO
#        echo "This is a new Repository that has been automatically created by forjj (POC)" > README.md
#        git add README.md
#        git commit -m 'initial commit'
#     fi
#     cd /devops/$REPO
#     if [ ! -d /devops/$REPO/.git ]
#     then
#        git init .
#        for FILE in README.md CONTRIBUTION.md
#        do
#          echo "Adding $FILE"
#          [ -f $FILE ] && git add FILE
#        done
#        git commit -m "Initial commit"
#     fi
#     RES="$(git remote -v | jq -Rs 'split("\n") | map(match("(.*)\t(.*) .*") | { "key": .captures[0].string, "value": .captures[1].string}) | from_entries')"
#
#     # Testing origin
#     if [ "$(json-get "$RES" '.origin')" != ""]
#     then
#        git remote remove origin        
#     fi
#     git remote add origin $FORK_SCP_LIKE
#     git push -u origin master
#     if [ "$(json-get "$RES" '.upstream')" != ""]
#     then
#        git remote remove upstream        
#     fi
#     echo "github creation done"
#     exit 0
#   fi
# 
#   if [ ! -d /devops/$REPO ]
#   then
#      cd /devops
#      git clone $FORK_SCP_LIKE $REPO
#      echo "github creation done"
#      exit 0
#   fi
#   cd /devops/$REPO
#   if [ ! -d .git ]
#   then
#      git init .
#   fi
#   RES="$(git remote -v | jq -Rs 'split("\n") | map(match("(.*)\t(.*) .*") | { "key": .captures[0].string, "value": .captures[1].string}) | from_entries')"
#   CUR_CONN="$(json-get "$RES" '.origin')"
#   case "$CUR_CONN" in
#      "$FORK_SSH_LIKE" | "$FORK_SCP_LIKE" | "$FORK_HTTP_LIKE")
#        git pull
#        echo "github creation done"
#        exit 0
#        ;;
#      null)
#        git remote add origin $FORK_SCP_LIKE
#        git fetch origin
#        echo "github creation done"
#        exit 0
#        ;;
#   esac    
#   echo "Existing local $REPO directory/clone may be in conflict with the existing remote repository. Please check it to be a git clone from $FORK_SCP_LIKE"
#   exit 1
#fi   
#}

function organization_found
{ 
 RES="$(github-get /user/orgs)"
 json-test "$RES" 'map(.login | test("'"$ORGA"'")) | any'
 return $?
}

function check_service
{
 ! organization_found $1 && return 1
}

################### MAIN

ACTION=$1
shift

INFRA_REPOS=false

parse_options "$@"

if [ "${GITHUB[TOKEN]}" != "" ]
then
   HEADER="Authorization: token ${GITHUB[TOKEN]}"
fi

if [ "${GITHUB[SERVER]}" = "" ]
then
   echo "Using public Github API"
   GITHUB[api_url]="https://api.github.com"
else
  GITHUB[api_url]="https://${GITHUB[SERVER]}/api/v3"
fi

RES="$(github-get /user)"
if [ $? -ne 200 ]
then
   echo "Unable to get current user information. Exiting."
   exit 1
fi

GITHUB[login]="$(json-get "$RES" '.login')"

if [ "${GITHUB[ORGANIZATION]}" != "" ]
then
   ORGA="${GITHUB[ORGANIZATION]}"
fi

case "$ACTION" in
  check)
    if ! check_service
    then
       printf "$FOUT" 404 "Organization $ORGA not found"
       exit 0
    else
       printf "$SOUT" "" "" 200 "Service up and running, repository found"
    fi
    ;;
  create)
    create_service
    ;;
  update)
    ;;
  maintain)
    ;;
  *)
    echo "Unknown action '$ACTION'"
    exit 1
    ;;
esac

