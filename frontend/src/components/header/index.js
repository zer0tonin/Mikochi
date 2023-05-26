import { h } from "preact";
import style from "./style.css";

import Search from "../search/search";

// TODO logo

const Header = ({ searchQuery, setSearchQuery }) => (
  <header class={style.header}>
    <nav role="navigation" aria-label="main navigation">
      <a
        href="/"
        onClick={() => setSearchQuery("")}
      >
        <img src="../../assets/logo.png" width="56" height="56" />
        Mikochi
      </a>
      <Search searchQuery={searchQuery} setSearchQuery={setSearchQuery} />
    </nav>
  </header>
);

export default Header;
