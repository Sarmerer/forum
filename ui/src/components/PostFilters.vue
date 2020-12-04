<template>
  <div>
    <transition name="fade" mode="out-in">
      <div v-if="!search" key="filter">
        <b-button @click="search = true" variant="dark">
          <b-icon-search></b-icon-search>
        </b-button>
        <b-button
          variant="dark"
          @click="sortCallback(), toast('b-toaster-bottom-center', true)"
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
            :disabled="sorter.throttled"
            @click="orderCallback(filter.orderBy)"
            v-b-tooltip.hover
            :title="`${sorter.asc ? 'Most' : 'Least'} ${filter.title}`"
            :variant="sorter.orderBy == filter.orderBy ? 'info' : 'dark'"
          >
            <b-icon :icon="filter.icon"></b-icon>
          </b-button>
        </b-button-group>
      </div>
      <div v-else key="search">
        <b-input-group v-if="search">
          <b-input-group-prepend>
            <b-button @click="search = false" variant="dark">
              <b-icon-arrow-left></b-icon-arrow-left>
            </b-button>
          </b-input-group-prepend>
          <b-form-input></b-form-input>
          <b-input-group-append>
            <b-button variant="outline-info" v-b-popover.top="'Popover!'">
              <b-icon-info></b-icon-info>
            </b-button>
          </b-input-group-append>
        </b-input-group>
      </div>
    </transition>
  </div>
</template>
<script>
export default {
  props: {
    orderCallback: {
      type: Function,
      required: true,
    },
    sortCallback: {
      type: Function,
      required: true,
    },
    sorter: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      search: false,
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
  methods: {
    toast(toaster, append = true) {
      var message = this.sorter.asc ? "ascending" : "descending";
      this.$bvToast.toast(`Posts sorted in ${message} order.`, {
        toaster: toaster,
        solid: true,
        appendToast: append,
        noCloseButton: true,
      });
    },
  },
};
</script>
<style lang="scss" scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter,
.fade-leave-to {
  opacity: 0;
}
</style>
