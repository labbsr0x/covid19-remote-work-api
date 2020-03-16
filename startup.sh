#!/bin/sh

./covid19-remote-work-api \
    --storageAPIURL=$STORAGE_API_URL \
    --certificateIssuerAPIURL=$CERTIFICATE_ISSUER_API_URL \
    --userRolesAPIURL=$USER_ROLED_API_URL \
    --conductorAPIURL=$CONDUCTOR_API_URL \
    --baseURL=$BASE_URL