import { h } from 'preact';
import { Link } from 'preact-router/match';
import style from './style.css';

// TODO logo

const Header = () => (
	<header class={style.header}>
		<nav role="navigation" aria-label="main navigation">
			<Link href="/">
				Mikochi
			</Link>
		</nav>
	</header>
);

export default Header;
