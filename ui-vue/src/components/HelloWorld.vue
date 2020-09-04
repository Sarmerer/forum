<template>
  <div>
    <h5>{{ msg }}</h5>
  </div>
</template>
<script>
export default {
  props: {
    msg: String,
  },
  data: function() {
    return {
      counter: this.msg,
    };
  },
};
</script>
