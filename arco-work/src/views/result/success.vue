<template>
  <div class="container">
    <Breadcrumb :items="['menu.result', 'menu.result.success']"/>
    <div class="wrapper">
      <a-result
          :subtitle="'表单提交成功！'"
          :title="'提交成功'"
          class="result"
          status="success"
      >


        <template #extra>
          <a-space :size="16" class="operation-wrap">
            <a-button key="again" type="secondary">
              {{ '打印结果' }}
            </a-button>
            <a-button key="back" type="primary">
              {{ '返回项目列表' }}
            </a-button>
          </a-space>
        </template>
      </a-result>

      <div class="steps-wrapper">
        <a-typography-paragraph bold>Processing</a-typography-paragraph>
        <a-steps :current="2" type="dot">
          <a-step
              :title="'提交申请'"
              description="2020/10/10 14:00:39"
          />
          <a-step
              :description="'进行中'"
              :title="'直属领导审核'"
          />
          <a-step
              :description="'未开始'"
              :title="'购买证书'"
          />
          <a-step
              :description="'未开始'"
              :title="'安全测试'"
          />
          <a-step
              :description="'未开始'"
              :title="'正式上线'"
          />
        </a-steps>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup></script>

<script lang="ts">
export default {
  name: 'Success',
};
</script>

<style lang="less" scoped>
.container {
  padding: 0 20px 20px 20px;
}

.wrapper {
  padding: 24px 150px;
  background-color: var(--color-bg-2);
  border-radius: 4px;
}

.result {
  margin: 150px 0 0 0;
}

.operation-wrap {
  margin-bottom: 40px;
  text-align: center;
}

.steps-wrapper {
  width: 100%;
  min-width: fit-content;
  margin-bottom: 150px;
  padding: 20px;
  background-color: rgb(var(--gray-1));
}
</style>

<style lang="less" scoped>
.mobile {
  .wrapper {
    padding: 24px 10px;

    .steps-wrapper {
      transform: scale(0.8);
    }
  }
}
</style>