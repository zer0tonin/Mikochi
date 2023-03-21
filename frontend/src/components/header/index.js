import { h } from 'preact';
import { Link } from 'preact-router/match';
import style from './style.css';

import Search from '../search/search';

// TODO logo

const Header = () => (
	<header class={style.header}>
		<nav role="navigation" aria-label="main navigation">
			<Link href="/">
				Mikochi
			</Link>
			<Search />
		</nav>
	</header>
);

export default Header;
