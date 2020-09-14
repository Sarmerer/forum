<template>
  <div>
    <div v-if="categories">
      <p v-for="c in categories" :key="c.ID" :id="c.ID">{{ c.Name }} | Posts: {{ c.UseCount }}</p>
    </div>
    <div v-else>
      <p>None</p>
    </div>
  </div>
</template>
<script>
import axios from "axios";
export default {
  data: function() {
    return {
      categories: [],
    };
  },
  created() {
    axios.get("categories").then((response) => (this.categories = response.data.data));
  },
  methods: {},
};
</script>
