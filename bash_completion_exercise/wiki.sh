!#/bin/bash

TITLE=$1

echo "Hello, you've started a min wiki scraper that will return the first line of a wikipedia page and list its subsections."

# curl -s "https://en.wikipedia.org/w/api.php?format=json&action=query&prop=extracts&explaintext=1&titles=$TITLE" | jq '.query.pages[].extract' | cut -d '.' -f -1 | cut -c 2-
# curl -s "https://en.wikipedia.org/wiki/$TITLE" 2>&1 | grep toclevel-1 | cut -d '#' -f 2- | cut -d '"' -f -1

function scrape_wiki {
	curl -s "https://en.wikipedia.org/w/api.php?format=json&action=query&prop=extracts&explaintext=1&titles=$1" | jq '.query.pages[].extract' | cut -d '.' -f -1 | cut -c 2-
	curl -s "https://en.wikipedia.org/wiki/$1" 2>&1 | grep toclevel-1 | cut -d '#' -f 2- | cut -d '"' -f -1
}

echo 'trying extra args'
# for ARG in $@; do
	# echo $ARG
	# scrape_wiki $ARG
# done
 


# psuedocode
# determine level at which to get info
    # get page of first argument
    # if no more arguments, get that sentence
    # else if more arguments, get sections of that page and search for section that matches last argument
        # if cannot match last argument, then move up one section
# get first sentence
# get sections




echo
echo 'next iteration'
echo "command arguments are: $@"
echo "number of command argumens is $#"

if [ $# == 1 ]; then
    echo 'received one argument'
    scrape_wiki $1
else
    echo "will give information on subsections because number of subsections is $#, greater than one"
    # NUMBER_SUBSECTIONS="$# - 1"
    # NUMBER_OF_ARGS=$#
    # echo "number of subsections is $NUMBER_SUBSECTIONS"
    # eval LAST_ARG=\$$#
    # echo "last argument is $LAST_ARG"
    # now get the sections associated with the last argument
    echo "second argument is $2"
    echo "getting toclevel from 'number' in API request"
    # curl -s "https://en.wikipedia.org/w/api.php?action=parse&page=$1&prop=sections&format=json" | jq '.parse.sections'
    TOC_LEVEL=$(curl -s "https://en.wikipedia.org/w/api.php?action=parse&page=Walrus&prop=sections&format=json" | jq -c '.parse.sections[] | select(.line | contains("'$2'")).number')
    echo "table of contents level is $TOC_LEVEL"
    echo "now using sed to get rid of extra quotes"
    TOC_LEVEL_CLEAN=$(echo "$TOC_LEVEL" | sed -E 's/"([^"]*)"/\1/')
    echo "cleaned valus is $TOC_LEVEL_CLEAN"
    echo "now getting all subsections"
    curl "https://en.wikipedia.org/wiki/$1" 2>&1 | grep toclevel- | grep -E "tocnumber\">$TOC_LEVEL_CLEAN\.+"
    echo "now cleaning final subsections"
    curl "https://en.wikipedia.org/wiki/$1" 2>&1 | grep toclevel- | grep -E "tocnumber\">$TOC_LEVEL_CLEAN\.+" | cut -d '#' -f 2- | cut -d '"' -f -1

    echo
    echo
    echo "now let's get the first part of the section"
    curl -s "https://en.wikipedia.org/w/api.php?action=parse&page=$1&section=1&prop=text&format=json" | jq -c '.parse.sections[] | select(.line | contains("'$2'")).number')
    # curl -s "https://en.wikipedia.org/w/api.php?action=parse&page=Walrus&section=1&prop=wikitext&format=json" | jq '.parse.wikitext."*"'
    # this is what I think should work but sed is giving me errors.
    # used this regex tool: https://regex101.com/
    # curl -s "https://en.wikipedia.org/w/api.php?action=parse&page=Walrus&section=1&prop=wikitext&format=json" | jq '.parse.wikitext."*"' | sed -E 's/   \\n[A-Z][A-Za-z\s'\[\]\.:]+\./   \1/  g'



    echo
    echo
    echo
    echo "so the final version of your work is..."
    curl "https://en.wikipedia.org/wiki/Walrus" 2>&1 | grep toclevel- | grep -E "tocnumber\">$TOC_LEVEL_CLEAN\.+" | cut -d '#' -f 2- | cut -d '"' -f -1
fi

# echo "sections call is"
# curl -s "https://en.wikipedia.org/w/api.php?action=parse&page=Walrus&prop=sections"