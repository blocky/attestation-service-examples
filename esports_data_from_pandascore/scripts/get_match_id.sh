#!/bin/bash

LEAGUE_NAME="starcraft-2-pl-invitational"
YEAR="2025"

LEAGUE_ID=$(curl -gs --request GET \
	--url "https://api.pandascore.co/leagues?filter[slug]=$LEAGUE_NAME" \
	--header 'accept: application/json' \
	--header "Authorization: Bearer $PANDASCORE_API_KEY" |
	jq ".[0].series[0].id")

TOURNAMENT_ID=$(curl -gs --request GET \
	--url "https://api.pandascore.co/series?filter[id]=$LEAGUE_ID" \
	--header 'accept: application/json' \
	--header "Authorization: Bearer $PANDASCORE_API_KEY" |
	jq ".[0].tournaments[] | select(.slug | contains(\"$YEAR\")) | .id")

MATCH_ID=$(curl -gs --request GET \
	--url "https://api.pandascore.co/tournaments?filter[id]=$TOURNAMENT_ID" \
	--header 'accept: application/json' \
	--header "Authorization: Bearer $PANDASCORE_API_KEY" |
	jq '.[0].matches[] | select(.name | contains("Grand final")) | .id')

echo "$MATCH_ID"
