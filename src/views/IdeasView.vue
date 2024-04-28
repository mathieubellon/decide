<template>
    <div>
        <router-link to="/">Home</router-link>
        <table>
            <thead>
                <tr>
                    <th>Title</th>
                    <th>Description</th>
                    <th>Workspace</th>
                    <th>priority</th>
                    <th>reach</th>
                    <th>votes</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="idea in ideas.ideas" :key="idea.id">
                    <td>{{ idea.title }}</td>
                    <td>{{ idea.description }}</td>
                    <td>{{ idea.workspace_id }}</td>
                    <td>{{ idea.priority }}</td>
                    <td>{{ idea.reach }}</td>
                    <td>{{ idea.votes }}</td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
const ideas = ref([]);

const fetchIdeas = async () => {
    try {
        const response = await fetch('/api/v1/ideas');
        const data = await response.json();
        ideas.value = data;
    } catch (error) {
        console.error('Error fetching ideas:', error);
    }
};

onMounted(fetchIdeas);

</script>