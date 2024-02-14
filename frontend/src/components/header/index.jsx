import { h } from "preact";
import "./style.css";

import Search from "../search";

const Header = ({ searchQuery, setSearchQuery }) => (
  <header class="header">
    <nav role="navigation" aria-label="main navigation">
      <a href="/" onClick={() => setSearchQuery("")}>
        <img src="/logo.png" width="56" height="56" />
        Mikochi
      </a>
      <Search searchQuery={searchQuery} setSearchQuery={setSearchQuery} />
    </nav>
  </header>
);

export default Header;
