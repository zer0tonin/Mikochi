import { h } from "preact";
import style from "./style.css";

const Search = ({ searchQuery, setSearchQuery }) => {
  const onSearchInput = (e) => {
    setSearchQuery(e.target.value);
  };

  const onSubmit = (e) => {
    e.preventDefault();
  };

  return (
    <form class={style.container} onSubmit={onSubmit}>
      <input
        class={style.searchInput}
        type="text"
        value={searchQuery}
        onInput={onSearchInput}
        aria-label="search"
      />
      <i class={`${style.searchIcon} gg-search`} />
    </form>
  );
};

export default Search;
