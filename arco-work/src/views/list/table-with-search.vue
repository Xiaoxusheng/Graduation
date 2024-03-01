<template>
  <div class="main-container">
    <TableBody ref="tableBody">
      <template #header>
        <TableHeader
            :show-filter="false"
            title="加班申请"
            @search="onSearch"
            @reset-search="onResetSearch"
        >
          <template #search-content>
            <a-form layout="inline" :model="{}">
              <a-form-item v-for="item of conditionItems" :key="item.key" :label="item.label">
                <template v-if="item.render">
                  <FormRender :render="item.render" :formItem="item"/>
                </template>
                <template v-else>
                  <template v-if="item.type === 'input'">
                    <a-input v-model="item.value.value" :placeholder="item.placeholder"/>
                  </template>
                  <template v-if="item.type === 'select'">
                    <a-select
                        v-model="item.value.value"
                        style="width: 150px"
                        :placeholder="item.placeholder"
                    >
                      <a-option
                          v-for="optionItem of item.optionItems"
                          :key="optionItem.value"
                          :value="optionItem.value"
                      >
                        {{ optionItem.label }}
                      </a-option>
                    </a-select>
                  </template>
                  <template v-if="item.type === 'date'">
                    <a-date-picker v-model="item.value.value"/>
                  </template>
                  <template v-if="item.type === 'time'">
                    <a-time-picker v-model="item.value.value" value-format="HH:mm:ss"/>
                  </template>
                  <template v-if="item.type === 'check-group'">
                    <a-checkbox-group v-model="item.value.value">
                      <a-checkbox v-for="it of item.optionItems" :value="it.value" :key="it.value">
                        {{ item.label }}
                      </a-checkbox>
                    </a-checkbox-group>
                  </template>
                </template>
              </a-form-item>
            </a-form>
          </template>
        </TableHeader>
      </template>
      <template #default>
        <a-table
            :bordered="{ wrapper: true, cell: true }"
            :row-selection="{ selectedRowKeys, showCheckedAll }"
            :loading="tableLoading"
            :data="dataList"
            :columns="tableColumns"
            :pagination="false"
            :rowKey="rowKey"
            @selection-change="onSelectionChange"
        >
          <template #columns>
            <a-table-column
                v-for="item of tableColumns"
                :key="item.key"
                :align="item.align"
                :title="(item.title as string)"
                :width="item.width"
                :data-index="(item.key as string)"
                :fixed="item.fixed"
            >
              <template v-if="item.key === 'index'" #cell="{ rowIndex }">
                {{ rowIndex + 1 }}
              </template>
              <template v-else-if="item.key === 'department_id'" #cell="{ record }">
                <a-tag :color="record.department_id === 1 ? 'green' : 'blue'">
                  {{ storedMap.get(record.department_id as number) }}
                </a-tag>
              </template>
              <template v-else-if="item.key === 'avatar'" #cell="{ record }">
                <a-avatar
                    :autocapitalize="30"
                    :style="{ backgroundColor: 'var(--color-primary-light-1)' }"
                >
                  {{ record.name }}
                </a-avatar>
              </template>
              <template v-else-if="item.key === 'status'" #cell="{ record }">
                <a-tag color="purple" size="small">请假</a-tag>
              </template>
              <template v-else-if="item.key === 'pass'" #cell="{ record }">
                <a-tag v-if="record.pass === 1" color="blue" size="small">通过</a-tag>
                <a-tag v-else-if="record.pass === 2" color="red" size="small">未通过</a-tag>
                <a-tag v-else color="purple" size="small">未审核</a-tag>
              </template>

              <template v-else-if="item.key === 'sex'" #cell="{ record }">
                <a-tag v-if="record.sex === 1" color="blue" size="small">男</a-tag>
                <a-tag v-else color="blue" size="small">女</a-tag>
              </template>
              <template v-else-if="item.key === 'actions' " #cell="{ record }">
                <a-space>
                  <a-button v-if="record.is_examine!==1" size="mini" type="primary" @click="onUpdateItem(record)">
                    审核
                  </a-button>
                  <a-button v-else disabled size="mini" type="primary" @click="onUpdateItem(record)">
                    已审核
                  </a-button>
                </a-space>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </template>
      <template #footer>
        <TableFooter :pagination="pagination"/>
      </template>
    </TableBody>
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
              <a-input v-model="item.value.value" :placeholder="item.placeholder" disabled></a-input>
            </template>
            <template v-if="item.type === 'textarea'">
              <a-textarea
                  v-model="item.value.value"
                  :auto-size="{ minRows: 2, maxRows: 5 }"
                  :placeholder="item.placeholder"
                  disabled
              />
            </template>
            <template v-if="item.type === 'switch'">
              <a-switch v-model="item.value.value" checked-children="开" un-checked-children="关"/>
            </template>
            <template v-if="item.type === 'date-range'">
              <a-range-picker v-model="item.value.value" :defaultValue="['2019-08-08 00:00:00', '2019-08-18 00:00:00']" :position="'tr'" disabled
                              showTime/>
            </template>
          </a-form-item>
        </a-form>
      </template>
    </ModalDialog>
  </div>
</template>

<script lang="ts">
import {get, post} from '@/api/http'
import {getLeaveApplicationList, leaveApplication,} from '@/api/url'
import {usePagination, useRowKey, useRowSelection, useTable, useTableColumn,} from '@/hooks/table'
import FormRender from '@/components/FormRender'
import {FormItem, ModalDialogType} from '@/types/components'
import {Form, Input, Message} from '@arco-design/web-vue'
import {defineComponent, h, onMounted, ref} from 'vue'
import type {Dayjs} from 'dayjs'
import useUserStore from "@/store/modules/user";
import ModalDialog from "@/components/ModalDialog.vue";

const conditionItems: Array<FormItem> = [
  {
    key: 'name',
    label: '用户姓名',
    type: 'input',
    placeholder: '请输入用户姓名',
    value: ref(''),
    reset: function () {
      this.value.value = ''
    },
    render: (formItem: FormItem) => {
      return h(Input, {
        placeholder: '请输入姓名',
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
  {
    key: 'sex',
    label: '用户姓别',
    value: ref(),
    type: 'select',
    placeholder: '请选择用户姓别',
    optionItems: [
      {
        label: '男',
        value: 1,
      },
      {
        label: '女',
        value: 2,
      },
    ],
    reset: function () {
      this.value.value = undefined
    },
  },
  {
    key: 'time',
    label: '创建时间',
    type: 'time',
    value: ref<string>(''),
  },
]
export default defineComponent({
  name: 'TableWithSearch',
  components: {
    ModalDialog,
    FormRender,
  },
  setup() {
    const searchForm = ref()
    const pagination = usePagination(doRefresh)
    const {selectedRowKeys, onSelectionChange, showCheckedAll} = useRowSelection()
    const table = useTable()
    const rowKey = useRowKey('id')
    const tableColumns = useTableColumn([
      table.indexColumn,
      {
        title: '员工姓名',
        key: 'name',
        dataIndex: 'name',
      },
      {
        title: '头像',
        key: 'avatar',
        dataIndex: 'avatar',
      },
      {
        title: '工号',
        key: 'uid',
        dataIndex: 'uid',
      },
      {
        title: '性别',
        key: 'sex',
        dataIndex: 'sex',
      },
      {
        title: '部门',
        key: 'department_id',
        dataIndex: 'department_id',
      },
      {
        title: '开始时间',
        key: 'start_time',
        dataIndex: 'start_time',
      },
      {
        title: '结束时间',
        key: 'end_time',
        dataIndex: 'end_time',
      },
      {
        title: '类型',
        key: 'status',
        dataIndex: 'status',
      },
      {
        title: '审核状态',
        key: 'pass',
        dataIndex: 'pass',
      },
      {
        title: '说明',
        key: 'reason',
        dataIndex: 'reason',
      },
      {
        title: '操作',
        key: 'actions',
        dataIndex: 'actions',
      }
    ])
    const userStore = useUserStore()
    const actionTitle = ref('请假审核')
    const modalDialogRef = ref<ModalDialogType | null>(null)
    const formItems = [
      {
        label: '姓名',
        key: 'name',
        type: 'input',
        placeholder: '请选择会议类型',
        value: ref(null),
      },
      {
        label: '工号',
        key: 'uid',
        type: 'input',
        placeholder: '请选择会议类型',
        value: ref(undefined),
      },
      {
        label: '性别',
        key: 'sex',
        type: 'input',
        placeholder: '请输入会议内容',
        value: ref(null),
      },
      {
        label: '起止日期',
        key: 'startEndDate',
        type: 'date-range',
        value: ref<Dayjs[]>([]),
        reset: function () {
          this.value.value = []
        },
      },
      {
        label: '请假原因',
        key: 'reason',
        placeholder: '请输入会议备注',
        type: 'textarea',
        value: ref(null),
      },
      {
        label: '是否通过',
        key: 'pass',
        type: 'switch',
        value: ref(false),
        required: true,
        placeholder: '请审核',
        validator: function () {
          if (!this.value.value) {
            Message.error(this.placeholder || '')
            return false
          }
          return true
        },
      },
    ] as FormItem[]
    const submitLoading = ref(false)
    const storedMapString = localStorage.getItem('departmentMap');
    const storedMapArray = JSON.parse(storedMapString);
    const storedMap = new Map(storedMapArray)

    let formRef = ref<typeof Form>()

    function doRefresh() {
      get({
        url: getLeaveApplicationList,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: () => {
          return {
            offset: pagination.page,
            limit: pagination.pageSize,
          }
        },
      }).then((res) => {
        res.data.forEach((i: any) => {
          const date = new Date(i.start_time);
          const year = date.getFullYear();
          const month = String(date.getMonth() + 1).padStart(2, "0");
          const day = String(date.getDate()).padStart(2, "0");
          const date1 = new Date(i.end_time);
          const year1 = date.getFullYear();
          const month1 = String(date.getMonth() + 1).padStart(2, "0");
          const day1 = String(date.getDate()).padStart(2, "0");
          i.end_time = `${year1}-${month1}-${day1}  ${date1.getHours() > 10 ? date1.getHours() : '0' + date1.getHours()}:${date1.getMinutes() > 10 ? date1.getMinutes() : '0' + date1.getMinutes()}:${date1.getSeconds() > 10 ? date1.getSeconds() : '0' + date1.getSeconds()}`
          i.start_time = `${year}-${month}-${day}  ${date.getHours() > 10 ? date.getHours() : '0' + date.getHours()}:${date.getMinutes() > 10 ? date.getMinutes() : '0' + date.getMinutes()}:${date.getSeconds() > 10 ? date.getSeconds() : '0' + date.getSeconds()}`
          return
        })
        table.handleSuccess(res)
        pagination.setTotalSize(res.data.length || 10)
      }).catch(error => console.log(error()))
    }

    function onSearch() {
      let data: any = conditionItems.reduce((pre, cur) => {
        ;(pre as any)[cur.key] = cur.value.value
        return pre
      }, {})
      const tableList = table.dataList.filter(i => {
        if (i.name === data.name || i.sex === data.sex) {
          return i
        }
      })
      table.handleSuccess(tableList)
      pagination.setTotalSize(tableList.length || 10)
    }

    function onResetSearch() {
      conditionItems.forEach((it) => {
        it.reset ? it.reset() : (it.value.value = '')
      })
    }


    // 审核
    function onUpdateItem(record: any) {
      // 处理数据
      /* [ "2023-02-08 21:17:58", "2024-03-20 21:14:58" ] */
      formItems.forEach(i => {
        if (i.key == 'name') {
          i.value.value = record.name
        }
        if (i.key == 'sex') {
          i.value.value = record.sex
        }
        if (i.key == 'uid') {
          i.value.value = record.uid
        }
        if (i.key == 'pass') {
          i.value.value = record.pass
        }
        if (i.key == 'reason') {
          i.value.value = record.reason
        }
        if (i.key == 'startEndDate') {
          i.value.value = [record.start_time, record.end_time]
        }
      })
      modalDialogRef.value?.toggle()
    }

    //   弹窗
    function onDataFormConfirm() {
      //
      if (formItems.every((it) => (it.validator ? it.validator() : true))) {
        let uid: number
        let pass: number
        formItems.forEach(i => {
          if (i.key == 'uid') {
            uid = i.value.value
          }
          if (i.key == "pass") {
            pass = i.value.value ? 1 : 2
          }
        })
        post({
          url: leaveApplication,
          headers: {
            Authorization: "Bearer " + userStore.token
          },
          data: () => {
            return {
              uid: uid,
              pass: pass as number,
            }
          },
        }).then((res) => {
          Message.success('审核成功')
          console.log(res)
          doRefresh()
        }).catch(error => {
          console.log(error)
          Message.success(error.toString(),)
        })

        modalDialogRef.value?.toggle()
      }

    }


    onMounted(doRefresh)
    return {
      ...table,
      rowKey,
      pagination,
      searchForm,
      tableColumns,
      conditionItems,
      onSearch,
      onResetSearch,
      selectedRowKeys,
      showCheckedAll,
      actionTitle,
      modalDialogRef,
      formItems,
      submitLoading,
      formRef,
      storedMap,
      onSelectionChange,
      onUpdateItem,
      onDataFormConfirm,
    }
  },
})
</script>

<style lang="less" scoped>
.form-wrapper {
  font-size: small;
}

.avatar-container {
  position: relative;
  width: 30px;
  height: 30px;
  margin: 0 auto;
  vertical-align: middle;

  .avatar {
    width: 100%;
    height: 100%;
    border-radius: 50%;
  }

  .avatar-vip {
    border: 2px solid #cece1e;
  }

  .vip {
    position: absolute;
    top: 0;
    right: -9px;
    width: 15px;
    transform: rotate(60deg);
  }
}

.gender-container {
  .gender-icon {
    width: 20px;
  }
}
</style>
