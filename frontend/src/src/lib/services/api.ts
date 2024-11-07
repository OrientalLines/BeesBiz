export async function getBeeCommunities() {
	const response = await fetch('/api/communities');
	return response.json();
}

export async function getHoneyHarvests() {
	const response = await fetch('/api/harvests');
	return response.json();
}
