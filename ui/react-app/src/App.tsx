import React from 'react';
import './App.css';

function App() {
	const [ item, setItem ] = React.useState([]);
	React.useEffect(() => {
		// POST request using fetch inside useEffect React hook
		const requestOptions = {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ title: 'React Hooks POST Request Example' })
		};
		fetch('https://localhost:4433/signin', requestOptions).then((response) => response.json()).then((data) => setItem(data));

		// empty dependency array means this effect will only run once (like componentDidMount in classes)
	}, []);
	return <h1>{item}</h1>;
}

export default App;
