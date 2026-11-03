import { h } from "preact";
import "./style.css";

import Search from "../search";
import Menu from "../menu";
import { useLocation } from "preact-iso";
import { getAssetUrl, getLocation } from "../../url";

const Header = ({ searchQuery, setSearchQuery }) => {
  const location = useLocation();
  const onHomeClick = () => {
    setSearchQuery("");
    location.route(getLocation("/"));
  };

  return (
    <header class="header">
      <nav role="navigation" aria-label="main navigation">
        <Menu onHomeClick={onHomeClick} />
        <a onClick={onHomeClick}>
          <img src={getAssetUrl("logo.png")} width="56" height="56" />
          Mikochi
        </a>
        <Search searchQuery={searchQuery} setSearchQuery={setSearchQuery} />
      </nav>
    </header>
  );
};

export default Header;
