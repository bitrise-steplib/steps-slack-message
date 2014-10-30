#!/bin/bash

THIS_SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "${THIS_SCRIPTDIR}/_utils.sh"
source "${THIS_SCRIPTDIR}/_formatted_output.sh"

# init / cleanup the formatted output
echo "" > "${formatted_output_file_path}"

if [ -z "${SLACK_CHANNEL}" ] ; then
	write_section_to_formatted_output "# Error"
	write_section_start_to_formatted_output '* Required input `$SLACK_CHANNEL` not provided!'
	exit 1
fi

if [ -z "${SLACK_FROM_NAME}" ] ; then
	write_section_to_formatted_output "# Error"
	write_section_start_to_formatted_output '* Required input `$SLACK_FROM_NAME` not provided!'
	exit 1
fi

if [ -z "${SLACK_MESSAGE_TEXT}" ] ; then
	write_section_to_formatted_output "# Error"
	write_section_start_to_formatted_output '* Required input `$SLACK_MESSAGE_TEXT` not provided!'
	exit 1
fi

if [ -z "${SLACK_WEBHOOK_URL}" ] ; then
	write_section_to_formatted_output "# Error"
	write_section_start_to_formatted_output '* Required input `$SLACK_WEBHOOK_URL` not provided!'
	exit 1
fi

write_section_to_formatted_output "# Configuration"
echo_string_to_formatted_output "* SLACK_CHANNEL: ${SLACK_CHANNEL}"
echo_string_to_formatted_output "* SLACK_FROM_NAME: ${SLACK_FROM_NAME}"
echo_string_to_formatted_output "* SLACK_MESSAGE_TEXT: ${SLACK_MESSAGE_TEXT}"
echo_string_to_formatted_output "* SLACK_WEBHOOK_URL: ${SLACK_WEBHOOK_URL}"

res=$(curl -s -X POST --data-urlencode "payload={\"channel\": \"${SLACK_CHANNEL}\", \"username\": \"${SLACK_FROM_NAME}\", \"text\": \"${SLACK_MESSAGE_TEXT}\"}" ${SLACK_WEBHOOK_URL})
# curl_ret_code=$?
# echo "Curl returned: ${curl_ret_code}"
if [ "${res}" == "ok" ] ; then
	write_section_to_formatted_output "# Success"
	echo_string_to_formatted_output "Message successfully sent."
	exit 0
fi

write_section_to_formatted_output "# Error"
write_section_to_formatted_output "Sending the message failed with the following error:"
echo_string_to_formatted_output "    ${res}"
exit 1