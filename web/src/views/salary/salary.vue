<template>
  <div style="width: 100%; height:100%;"/>
  <TableBody>
    <template #header>
      <TableHeader ref="tableHeaderRef" :show-filter="false">
        <template #table-config :title="'工资列表'">
          <a-form :model="{}" layout="inline">
            <a-form-item v-for="item of conditionItems" :key="item.key" :label="item.label">
              <template v-if="item.render">
                <FormRender :formItem="item" :render="item.render"/>
              </template>
              <template v-else>
                <template v-if="item.type === 'date'">
                  <a-month-picker v-model="item.value.value"/>
                </template>
              </template>
            </a-form-item>
          </a-form>
          <a-space>
            <a-button size="small" type="primary" @click="getSalarys">查询</a-button>
            <a-button size="small" type="primary" @click="input">输入考勤</a-button>
          </a-space>
        </template>
      </TableHeader>
    </template>
    <template #default>
      <a-table
          :column-resizable="false"
          :data="dataList"
          :hoverable="true"
          :loading="tableLoading"
          :pagination="false"
          :rowKey="rowKey"
          :scroll="{ y: tableHeight }"
          :stripe="true"
          size="small"
          table-layout-fixe
          @selection-change="onSelectionChange"
      >
        <template #columns>
          <a-table-column
              v-for="item of tableColumns"
              :key="item.identity"
              :align="item.align"
              :data-index="(item.key as string)"
              :fixed="item.fixed"
              :title="(item.title as string)"
              :width="item.width"
          >
            <template v-if="item.key === 'index'" #cell="{ rowIndex }">
              {{ rowIndex + 1 }}
            </template>
            <template v-else-if="item.key === 'attendance_hours'" #cell="{ record }">
              <a-tag :color="'green'">
                {{ record.attendance_hours.toFixed(1) }}
              </a-tag>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </template>

    <template #footer>
      <TableFooter ref="tableFooterRef" :pagination="pagination"/>
    </template>

  </TableBody>
  <!--    继续-->
  <ModalDialog ref="modalDialogRef" :title="actionTitle" @confirm="onDataFormConfirm">
    <template #content>
      <a-form :label-col-props="{ span: 5 }" :model="formRef">
        <a-form-item
            v-for="item of formItems"
            :key="item.key"
            :label="item.label"
            :row-class="[item.required ? 'form-item__require' : 'form-item__no_require']"
            label-align="left"
        >
          <template v-if="item.type === 'input'">
            <a-input v-model.number="item.value.value" :placeholder="item.placeholder"></a-input>
          </template>
          <template v-if="item.type === 'input-number'">
            <a-input-number v-model="item.value.value" :placeholder="item.placeholder" mode="button"></a-input-number>
          </template>
        </a-form-item>
      </a-form>
    </template>
  </ModalDialog>

</template>

<script lang="ts">
import {getSalary, getSalaryList, salaryInfo,} from '@/api/url'
import {usePagination, useRowKey, useRowSelection, useTable, useTableColumn, useTableHeight,} from '@/hooks/table'
import {Form, Input, Message} from '@arco-design/web-vue'
import {defineComponent, getCurrentInstance, h, onMounted, ref} from 'vue'
import AddButton from "@/components/AddButton.vue";
import useUserStore from "@/store/modules/user";
import {FormItem, ModalDialogType} from "@/types/components";
import usePost from '@/hooks/usePost'
import useGet from "@/hooks/useGet";
import FormRender from "@/components/FormRender";
import type {Dayjs} from "dayjs";


export default defineComponent({
  components: {FormRender, AddButton},
  name: 'UserList',
  setup: function () {
    const table = useTable()
    const rowKey = useRowKey('id')
    const pagination = usePagination(doRefresh)
    const {selectedRowKeys, onSelectionChange} = useRowSelection()
    const conditionItems: Array<FormItem> = [
      {
        key: 'uid',
        label: '员工工号',
        type: 'input',
        placeholder: '请输入员工工号',
        value: ref(''),
        reset: function () {
          this.value.value = ''
        },
        render: (formItem: FormItem) => {
          return h(Input, {
            placeholder: '请输入员工工号',
            modelValue: formItem.value.value,
            'onUpdate:modelValue': (value) => {
              formItem.value.value = value
            },
          })
        },
      },
      {
        key: 'date',
        label: '创建日期',
        type: 'date',
        value: ref<Dayjs>(),
      },
    ]
    const formRef = ref<typeof Form>()
    const formModel = ref({})
    const modalDialogRef = ref<ModalDialogType | null>(null)
    const actionTitle = ref('添加员工考勤信息')
    const tableColumns = useTableColumn([
      table.indexColumn,
      {
        title: '工号',
        key: 'uid',
        dataIndex: 'uid',
      },
      {
        title: '月份',
        key: 'date',
        dataIndex: 'date',
      },
      {
        title: '早退迟到次数',
        key: 'count',
        dataIndex: 'count',
      },
      {
        title: '上班时长',
        key: 'attendance_hours',
        dataIndex: 'attendance_hours',
      },
      {
        title: '公司补贴',
        key: 'subsidy',
        dataIndex: 'subsidy',
      }, {
        title: '带薪休假时长',
        key: 'paid_leave',
        dataIndex: 'paid_leave',
      },
      {
        title: '其他',
        key: 'other',
        dataIndex: 'other',
      },
      {
        title: '工资',
        key: 'total',
        dataIndex: 'total',
      },
      {
        title: '补贴',
        key: 'subsidy',
        dataIndex: 'subsidy',
      },
      {
        title: '生成时间',
        key: 'CreatedAt',
        dataIndex: 'CreatedAt',
        width: 200,
      },
    ])
    const expandAllFlag = ref(true)
    const userStore = useUserStore()
    const post = usePost()
    const get = useGet()
    const formItems = [
      {
        label: '工号',
        key: 'uid',
        type: 'input',
        placeholder: '请输入工号',
        value: ref(undefined),
        required: true,
        validator: function () {
          if (!this.value.value) {
            Message.error(this.placeholder || '')
            return false
          }
          return true
        },
      },
      {
        label: '异常次数',
        key: 'count',
        type: 'input-number',
        placeholder: '选择次数',
        value: ref(null),
        required: true,
        validator: function () {
          if (!this.value.value) {
            Message.error(this.placeholder || '')
            return false
          }
          return true
        },
      },
      {
        label: '带薪休假',
        key: 'paidLeave',
        placeholder: '请输入带薪休假时长',
        type: 'input',
        value: ref(null),
        required: true,
        validator: function () {
          if (!this.value.value) {
            Message.error(this.placeholder || '')
            return false
          }
          return true
        },
      },
      {
        label: '工作时长',
        key: 'attendanceHours',
        placeholder: '请输入工作时长',
        type: 'input',
        value: ref(null),
        required: true,
        validator: function () {
          if (!this.value.value) {
            Message.error(this.placeholder || '')
            return false
          }
          return true
        },
      },
      {
        label: '补贴',
        key: 'other',
        placeholder: '请输入补贴',
        type: 'input',
        value: ref(null),
        required: true,
        validator: function () {
          if (!this.value.value) {
            Message.error(this.placeholder || '')
            return false
          }
          return true
        },
      }, {
        label: '社保',
        key: 'subsidy',
        placeholder: '请输入',
        type: 'input',
        value: ref(null),
        required: true,
        validator: function () {
          if (!this.value.value) {
            Message.error(this.placeholder || '')
            return false
          }
          return true
        },
      },

    ] as FormItem[]

    function doRefresh() {
      get({
        url: getSalaryList,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: () => {
          return {
            time: new Date().getFullYear() + '-' + (new Date().getMonth() + 1 < 10 ? '0' + (new Date().getMonth() + 1) : new Date().getMonth() + 1),
            offset: pagination.page,
            limit: pagination.pageSize,
          }
        },
      }).then((res) => {
        res.data.forEach((i: any) => {
          const date = new Date(i.CreatedAt);
          const year = date.getFullYear();
          const month = String(date.getMonth() + 1).padStart(2, "0");
          const day = String(date.getDate()).padStart(2, "0");
          i.CreatedAt = `${year}-${month}-${day}  ${date.getHours() > 10 ? date.getHours() : '0' + date.getHours()}:${date.getMinutes() > 10 ? date.getMinutes() : '0' + date.getMinutes()}:${date.getSeconds() > 10 ? date.getSeconds() : '0' + date.getSeconds()}`
        })
        table.handleSuccess(res)
        console.log(res)
        // pagination.setTotalSize(res.count)
      }).catch(error => {
        console.log(error)
      })
    }


    function getSalarys() {
      let uid: any
      if (conditionItems[0].value.value == "") {
        Message.error("值不能为空")
        return
      }
      conditionItems.forEach(i => {
        if (i.key == "uid") {
          uid = i.value.value
        }
      })
      get({
        url: getSalary,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: () => {
          return {
            uid: uid,
          }
        },
      }).then((res) => {
        res.data.forEach((i: any) => {
          const date = new Date(i.CreatedAt);
          const year = date.getFullYear();
          const month = String(date.getMonth() + 1).padStart(2, "0");
          const day = String(date.getDate()).padStart(2, "0");
          i.CreatedAt = `${year}-${month}-${day}  ${date.getHours() > 10 ? date.getHours() : '0' + date.getHours()}:${date.getMinutes() > 10 ? date.getMinutes() : '0' + date.getMinutes()}:${date.getSeconds() > 10 ? date.getSeconds() : '0' + date.getSeconds()}`
        })
        table.handleSuccess(res)
        pagination.setTotalSize(1)
      }).catch(error => {
        Message.error(error.toString())
        console.log(error)
      })
      conditionItems.forEach(i => {
        i.value.value = ''
      })
    }

    function input() {
      modalDialogRef.value?.toggle()
    }

    function onDataFormConfirm() {
      post({
        url: salaryInfo,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: formItems.reduce((pre, cur) => {
          ;(pre as any)[cur.key] = cur.value.value
          return pre
        }, {})
      }).then((res) => {
        Message.success("添加成功")
      }).catch(error => {
        Message.error(error.toString())
      })
      modalDialogRef.value?.toggle()
      formItems.forEach((i: any) => {
        i.value.value = ''
      })
    }

    onMounted(async () => {
      table.tableHeight.value = await useTableHeight(getCurrentInstance())
      doRefresh()
    })
    return {
      ...table,
      rowKey,
      selectedRowKeys,
      onSelectionChange,
      expandAllFlag,
      tableColumns,
      pagination,
      formModel,
      actionTitle,
      getSalarys,
      modalDialogRef,
      onDataFormConfirm,
      input,
      formRef,
      conditionItems,
      formItems,
    }
  },
})
</script>
