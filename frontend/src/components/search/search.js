import { h } from "preact";
import { useState } from "preact/hooks";
import style from "./style.css";

const Search = ({ searchQuery, setSearchQuery }) => {
  const [value, setValue] = useState(searchQuery);

  const onSearchInput = (e) => {
    setValue(e.target.value);

    // only trigger search when the user stops typing
    const timer = setTimeout(() => setSearchQuery(e.target.value), 500);
    return () => clearTimeout(timer);
  };

  const onSubmit = (e) => {
    e.preventDefault();
  };

  return (
    <form class={style.container} onSubmit={onSubmit}>
      <input
        class={style.searchInput}
        type="text"
        value={value}
        onInput={onSearchInput}
        aria-label="search"
      />
      <i class={`${style.searchIcon} gg-search`} />
    </form>
  );
};

export default Search;
