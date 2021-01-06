<template>
  <b-input-group size="sm">
    <template #prepend>
      <b-dropdown :text="`Find ${options.active}`" variant="dark" size="sm">
        <b-dropdown-item
          v-for="(filter, index) in options.filters"
          :key="index"
          variant="dark"
          :disabled="filter.name == options.active"
          @click="options.active = filter.name"
          >{{ filter.name }}</b-dropdown-item
        >
      </b-dropdown>
      <b-dropdown text="by author" variant="dark" size="sm">
        <b-dropdown-item variant="dark">ban</b-dropdown-item>
      </b-dropdown>
    </template>
    <b-form-input></b-form-input>

    <template #append>
      <b-button variant="dark"><b-icon-search></b-icon-search></b-button>
    </template>
  </b-input-group>
</template>
<script>
export default {
  name: "SearchBar",
  data() {
    return {
      options: {
        active: "post",
        filters: [{ name: "post" }, { name: "user" }],
      },
    };
  },
};
</script>
