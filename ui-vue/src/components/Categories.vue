<template>
  <div v-if="categories && categories.lenght > 0">
    <p v-for="c in categories" :key="c.ID" :id="c.ID">{{ c.Name }} | Posts: {{ c.UseCount }}</p>
  </div>
  <div v-else>
    <p>None</p>
  </div>
</template>
<script>
import axios from "axios";
export default {
  data: function() {
    return {
      categories: null,
    };
  },
  created() {
    axios.get("categories").then((response) => (this.categories = response.data.data));
  },
  methods: {},
};
</script>
