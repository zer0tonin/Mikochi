import { h } from 'preact';
import { Link } from 'preact-router/match';
import style from './style.css';

import Search from '../search/search';

// TODO logo

const Header = ({ searchQuery, setSearchQuery }) => (
	<header class={style.header}>
		<nav role="navigation" aria-label="main navigation">
			<Link href="/">
				Mikochi
			</Link>
			<Search searchQuery={searchQuery} setSearchQuery={setSearchQuery} />
		</nav>
	</header>
);

export default Header;
