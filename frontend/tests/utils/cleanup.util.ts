import axios from 'axios';

export async function cleanupBackend() {
	await axios.post('/api/test/reset');
}
