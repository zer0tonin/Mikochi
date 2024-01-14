import { h } from "preact";
import { useEffect, useState } from "preact/hooks";
import "./style.css";

const Search = ({ searchQuery, setSearchQuery }) => {
  const [value, setValue] = useState(searchQuery);

  useEffect(() => {
    if (searchQuery === "") {
      setValue(searchQuery);
    }
  }, [searchQuery]);

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
    <form class="container" onSubmit={onSubmit}>
      <input
        class="searchInput"
        type="text"
        value={value}
        onInput={onSearchInput}
        aria-label="search"
      />
      <i class="searchIcon gg-search" />
    </form>
  );
};

export default Search;
