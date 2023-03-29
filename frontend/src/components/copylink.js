import { h } from 'preact';

import Icon from './icon';

const CopyLink = ({ filePath }) => {
	const copyToClipboard = () => {
		navigator.clipboard.writeText(`${window.location.protocol}//${window.location.hostname}${window.location.port == "" ? "" : ":" + window.location.port}/api/stream${filePath}`)
	}
	return (
		<Icon name="copy" onClick={copyToClipboard} />
	);
};

export default CopyLink;
