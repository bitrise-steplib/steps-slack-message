#!/bin/bash

THIS_SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "${THIS_SCRIPTDIR}/_bash_utils/utils.sh"
source "${THIS_SCRIPTDIR}/_bash_utils/formatted_output.sh"

# init / cleanup the formatted output
echo "" > "${formatted_output_file_path}"

if [ -z "${channel}" ] ; then
	write_section_to_formatted_output '*Notice: `$channel` is not provided!*'
fi

if [ -z "${from_username}" ] ; then
	write_section_to_formatted_output '*Notice: `$from_username` is not provided!*'
fi

if [ -z "${message}" ] ; then
	write_section_to_formatted_output "# Error"
	write_section_start_to_formatted_output '* Required input `$message` not provided!'
	exit 1
fi

if [ -z "${webhook_url}" ] ; then
	write_section_to_formatted_output "# Error"
	write_section_start_to_formatted_output '* Required input `$webhook_url` not provided!'
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
