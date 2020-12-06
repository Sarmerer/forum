<template>
  <span>
    <span
      v-if="tooltip"
      v-b-tooltip="{
        title: datetime,

        trigger: 'hover',
        placement: tooltip || 'top',
      }"
    >
      {{ ago }}
    </span>
    <span v-else>
      {{ ago }}
    </span>
  </span>
</template>
<script>
import moment from "moment-shortformat";

export default {
  name: "TimeAgo",
  props: {
    datetime: { required: true },
    tooltip: String,
  },
  created() {
    this.ago = moment
      .unix(this.datetime)
      //.fromNow()
      .short(this.isMobile());
  },

  data() {
    return {
      ago: null,
    };
  },
};
</script>
