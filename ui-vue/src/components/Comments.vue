<template>
  <div class="card">
    <b-form @submit="leaveComment">
      <b-input-group class="mt-1">
        <b-textarea
          type="text"
          placeholder="Comment this post"
          v-model="form.comment"
          rows="1"
          max-rows="3"
          size="sm"
          no-resize
        ></b-textarea>
        <b-input-group-append>
          <b-button
            variant="outline-light"
            type="submit"
            :disabled="form.comment.length < 5 || form.comment.length > 200"
            >Say</b-button
          >
        </b-input-group-append>
      </b-input-group>
      <div v-if="form.comment.length > 0">
        <small v-if="form.comment.length >= 5 && form.comment.length <= 200"
          >{{ form.comment.length }}/200</small
        >
        <small v-else style="color: red">{{ form.comment.length }}/200</small>
      </div>
    </b-form>
    <div v-for="(comment, index) in comments" :key="index">
      <div style="margin: 0.3rem; position: relative">
        <p
          v-if="index != currComment.editing"
          style="margin-bottom: 0px; display: block; margin-right: 4rem"
        >
          {{ comment.author_name }} says: {{ comment.content }}
        </p>
        <b-form-textarea
          v-else
          v-model="currComment.editingContent"
          rows="1"
          no-resize
          max-rows="10"
          style="margin-bottom: 0px; display: block; width: 85%"
        ></b-form-textarea>
        <small v-b-tooltip.hover :title="comment.created"
          ><timeago :datetime="comment.created" :auto-update="60"></timeago
        ></small>
        <b-button-group
          v-if="
            user && (post.author_id == user.id || user.role > 0) && index != currComment.editing
          "
          size="sm"
          class="controls-button"
          style="position: absolute; right: 0px; top: 10px"
        >
          <b-button
            size="sm"
            lg="1"
            variant="light"
            class="controls-button"
            :disabled="currComment.deleting"
            @click="
              () => {
                currComment.editing = index;
                currComment.editingContent = comment.content;
              }
            "
          >
            <img src="@/assets/svg/post/edit.svg" alt="edit" srcset="" />
          </b-button>
          <b-button
            variant="danger"
            :disabled="currComment.deleting"
            class="controls-button"
            @click="deleteComment(comment.id, index)"
            ><img src="@/assets/svg/post/delete.svg" alt="delete" srcset=""
          /></b-button>
        </b-button-group>
        <b-button-group
          size="sm"
          vertical
          v-if="index == currComment.editing"
          style="position: absolute; right: 0px; top: 2px"
        >
          <b-button
            :disabled="currComment.editingContent == comment.content"
            variant="success"
            @click="updateComment(comment.id)"
          >
            Save
          </b-button>
          <b-button variant="outline-danger" @click="currComment.editing = -1">Cancel</b-button>
        </b-button-group>
      </div>
    </div>
  </div>
</template>
