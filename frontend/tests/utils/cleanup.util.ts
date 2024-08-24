import axios from 'axios';

export async function cleanupBackend() {
	await axios.post('http://localhost/api/test/reset');
}
