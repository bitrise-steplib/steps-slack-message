#!/bin/bash

THIS_SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "${THIS_SCRIPTDIR}/_bash_utils/utils.sh"
source "${THIS_SCRIPTDIR}/_bash_utils/formatted_output.sh"

# init / cleanup the formatted output
echo "" > "${formatted_output_file_path}"

if [ -z "${SLACK_CHANNEL}" ] ; then
	write_section_to_formatted_output '*Notice: `$SLACK_CHANNEL` is not provided!*'
fi

if [ -z "${SLACK_FROM_NAME}" ] ; then
	write_section_to_formatted_output '*Notice: `$SLACK_FROM_NAME` is not provided!*'
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

resp=$(go run "${THIS_SCRIPTDIR}/step.go")
ex_code=$?

if [ ${ex_code} -eq 0 ] ; then
	echo "${resp}"
	write_section_to_formatted_output "# Success"
	echo_string_to_formatted_output "Message successfully sent."
	exit 0
fi

write_section_to_formatted_output "# Error"
write_section_to_formatted_output "Sending the message failed with the following error:"
echo_string_to_formatted_output "${resp}"
exit 1
