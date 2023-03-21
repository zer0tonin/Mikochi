import { h } from 'preact';
import style from './style.css';

import Icon from '../icon'

const Search = () => {
	return (
	  <form class={style.container}>
		<input class={style.searchInput} type="text" />
		<i class={`${style.searchIcon} gg-search`} />
	  </form>
	)
}

export default Search
