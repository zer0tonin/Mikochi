import { h } from 'preact';

const iconStyle = {
	margin: "0.1em",
	display: "inline-block",
};

const Icon = ({ name }) => (
	<i style={iconStyle} class={`gg-${name}`}></i>
);

export default Icon;
