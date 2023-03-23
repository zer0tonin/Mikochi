import { h } from 'preact';

const Path = ({ fileInfo }) => {
	const splitPath = fileInfo.path.split('/')

	return (
		<>
			{
				splitPath
					.map((val, i) => {
						if (i == splitPath.length - 1 && !fileInfo.isDir) {
							return (<>{val}</>)
						}
						return (<a href={splitPath.slice(0, i + 1).join('/') + '/'}>{val}</a>)
					})
					.reduce((acc, val) => (
						acc === null ? [val] : [...acc, ' / ', val]
					), null)
			}
		</>
	);
};

export default Path;
