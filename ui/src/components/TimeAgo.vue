<template>
  <span>
    <span
      v-if="tooltip"
      v-b-tooltip="{
        title: dateFormatted,
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
    this.ago = moment.unix(this.datetime).short();
    this.dateFormatted = moment
      .unix(this.datetime)
      .format("ddd, D MMM y, hh:mm");
  },

  data() {
    return {
      ago: null,
      dateFormatted: null,
    };
  },
};
</script>
