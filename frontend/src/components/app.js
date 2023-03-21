import { h } from 'preact';
import { Router } from 'preact-router';

import Header from './header';

// Code-splitting is automated for `routes` directory
import Directory from '../routes/directory/directory';

const App = () => (
	<div id="app">
		<Header />
		<main>
			<Router>
				<Directory path="/" />
				<Directory path="/:dirPath" />
			</Router>
		</main>
	</div>
);

export default App;
