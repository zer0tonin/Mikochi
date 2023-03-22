import { h } from 'preact';
import { useState, useEffect } from 'preact/hooks';
import style from './style.css';

import Icon from '../icon'

const Search = ({ searchQuery, setSearchQuery }) => {
	const onSearchInput = (e) => {
		setSearchQuery(e.target.value)
	}

	const onSubmit = (e) => {
		e.preventDefault()
	}

	return (
	  <form class={style.container} onSubmit={onSubmit}>
		<input class={style.searchInput} type="text" value={searchQuery} onInput={onSearchInput} />
		<i class={`${style.searchIcon} gg-search`} />
	  </form>
	)
}

export default Search
