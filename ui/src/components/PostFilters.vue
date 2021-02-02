<template>
  <div>
    <b-button
      variant="dark"
      @click="$emit('sort-event')"
      :disabled="sorter.throttled"
      class="mx-2"
      v-b-tooltip.hover
      :title="sorter.asc ? 'Ascending' : 'Descending'"
    >
      <b-icon :icon="sorter.asc ? 'sort-up' : 'sort-down-alt'"> </b-icon>
    </b-button>
    <b-button-group>
      <b-button
        v-for="filter in filters"
        :key="filter.orderBy"
        :disabled="sorter.throttled || sorter.orderBy === filter.orderBy"
        @click="$emit('order-event', filter.orderBy)"
        v-b-tooltip.hover
        :title="`${sorter.asc ? 'Most' : 'Least'} ${filter.title}`"
        :variant="sorter.orderBy == filter.orderBy ? 'info' : 'dark'"
      >
        <b-icon :icon="filter.icon"></b-icon>
      </b-button>
    </b-button-group>
  </div>
</template>
<script>
export default {
  props: {
    sorter: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      filters: [
        { orderBy: "rating", title: "likes", icon: "heart" },
        { orderBy: "created", title: "recent", icon: "clock" },
        { orderBy: "comments_count", title: "comments", icon: "chat" },
        {
          orderBy: "total_participants",
          title: "participants",
          icon: "people",
        },
      ],
    };
  },
};
</script>
