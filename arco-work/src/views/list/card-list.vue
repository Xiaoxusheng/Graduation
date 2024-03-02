<template>
  <div class="main-container">
    <TableBody>
      <template #header>
        <TableHeader :show-filter="false">
          <template #date-content>
            <a-space class="ml-4">
              请选择考勤日期 ：
              <a-date-picker v-model="date" style="width: 200px;" @change="check"/>
            </a-space>
          </template>
          <template #table-config>
            <a-space>
              <SortableTable :columns="tableColumns" class="ml-4" @update="onUpdateTable"/>
              <TableConfig
                  @refresh="doRefresh"
                  @update-border="onUpdateBorder"
                  @update-striped="onUpdateStriped"
              />
            </a-space>
          </template>
        </TableHeader>
      </template>
      <template #default>
        <a-table
            :bordered="{ wrapper: bordered, cell: bordered }"
            :data="dataList"
            :loading="tableLoading"
            :pagination="false"
            :row-class-name="rowClassNameFun"
            :row-key="rowKey"
            :stripe="(striped as boolean)"
        >
          <template #columns>
            <a-table-column
                v-for="item of tableColumns"
                :key="item.key"
                :data-index="(item.key as string)"
                v-bind="item"
            >
              <template v-if="item.key === 'index'" #cell="{ rowIndex }">
                {{ rowIndex + 1 }}
              </template>
              <template v-else-if="item.key === 'sex'" #cell="{ record }">
                <a-tag v-if="record.sex === 1" color="blue" size="small">男</a-tag>
                <a-tag v-else color="blue" size="small">女</a-tag>
              </template>
              <template v-else-if="item.key === 'avatar'" #cell="{ record }">
                <a-avatar :size="30" :style="{ backgroundColor: 'var(--color-primary-light-1)' }">
                  {{ record.name }}
                </a-avatar>
              </template>
              <template v-else-if="item.key === 'department_id'" #cell="{ record }">
                <a-tag :color="record.department_id === 1 ? 'green' : 'blue'">
                  {{ storedMap.get(record.department_id as number) }}
                </a-tag>
              </template>
              <template v-else-if="item.key === 'status'" #cell="{ record }">
                <a-tag color="red" size="small">{{ map[record.status] }}</a-tag>
              </template>
              <template v-else-if="item.key === 'url'" #cell="{ record }">
                <a-image
                    :preview-props="{
                     actionsLayout: ['rotateRight', 'zoomIn', 'zoomOut'],
               }"
                    :src="record.url"
                    width="100"
                />
              </template>
              <template v-else-if="item.key === 'pass'" #cell="{ record }">
                <a-tag v-if="record.pass === 1" color="blue" size="small">通过</a-tag>
                <a-tag v-else-if="record.pass === 2" color="red" size="small">未通过</a-tag>
                <a-tag v-else color="purple" size="small">未审核</a-tag>
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
        <TableFooter :pagination="pagination" position="end"/>
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
              <a-switch v-model.number="item.value.value" checked-children="1" un-checked-children="2"/>
            </template>
            <template v-if="item.type === 'date-range'">
              <a-range-picker v-model="item.value.value" :defaultValue="['2019-08-08 00:00:00', '2019-08-18 00:00:00']"
                              :position="'tr'" disabled
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
import {getAllClockIn, overtimeApplication} from '@/api/url'
import {usePagination, useRowKey, useTable, useTableColumn} from '@/hooks/table'
import {FormItem, ModalDialogType, TablePropsType} from '@/types/components'
import {sortColumns} from '@/utils'
import {Form, Message, Modal} from '@arco-design/web-vue'
import {defineComponent, onMounted, reactive, ref} from 'vue'
import useUserStore from "@/store/modules/user";
import ModalDialog from "@/components/ModalDialog.vue";
import type {Dayjs} from "dayjs";

export default defineComponent({
  name: 'TableCustom',
  components: {ModalDialog},
  setup() {
    const table = useTable()
    const pagination = usePagination(doRefresh)
    const rowKey = useRowKey('id')
    const tableColumns = reactive(
        useTableColumn([
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
            title: '图片',
            key: 'url',
            dataIndex: 'url',
          },

        ])
    )
    const userStore = useUserStore()
    const actionTitle = ref('加班审核')
    const modalDialogRef = ref<ModalDialogType>(null)

    const storedMapString = localStorage.getItem('departmentMap');
    const storedMapArray = JSON.parse(storedMapString);
    const storedMap = new Map(storedMapArray)
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

    ] as FormItem[]
    const map = {
      1: "缺勤",
      2: "打卡",
      3: "迟到",
      4: "加班",
      5: "补卡",
      6: "请假",
    }

    const date = ref('')

    let formRef = ref<typeof Form>()
// 获取当前日期
    const today = new Date();

// 设置时间为 0 点
    today.setHours(0, 0, 0, 0);

    const time = Math.floor(today.getTime() / 1000)

    function doRefresh() {
      console.log(time)
      get({
        url: getAllClockIn,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: () => {
          return {
            time: time
          }
        },
      }).then((res) => {
        console.log(res)
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
      }).catch(error => {
        Message.error(error.toString())
      })
    }

    function onDeleteItem(item: any) {
      if (item) {
        Modal.confirm({
          content: '是否要删除此数据，删除后不恢复？',
          okText: '删除',
          onOk: () => {
            Message.success('模拟删除成功，参数为：' + item.id)
          },
        })
      }
    }

    function onUpdateTable(newColumns: Array<TablePropsType>) {
      sortColumns(tableColumns, newColumns)
    }

    function onUpdateBorder(isBordered: boolean) {
      table.bordered.value = isBordered
    }

    function onUpdateStriped(isStriped: boolean) {
      table.striped.value = isStriped
    }

    function rowClassNameFun(_record: any, index: number) {
      return index % 2 === 1 && table.striped.value ? 'table-striped' : null
    }

    function onUpdateItem(record: any) {
      // 处理数据
      console.log(record)

      /* [ "2023-02-08 21:17:58", "2024-03-20 21:14:58" ] */
      formItems.forEach(i => {
        if (i.key == 'name') {
          i.value.value = record.name as string
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
        url: overtimeApplication,
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

    // 时间
    function check(data: any) {
      const today = new Date(date.value);
      today.setHours(0, 0, 0, 0);
      let time = Math.floor(today.getTime() / 1000)
      get({
        url: getAllClockIn,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: () => {
          return {
            time: time
          }
        },
      }).then((res) => {
        console.log(res)
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
      }).catch(error => {
        Message.error(error.toString())
      })
    }
    onMounted(doRefresh)
    return {
      ...table,
      rowKey,
      tableColumns,
      pagination,
      storedMap,
      actionTitle,
      formRef,
      formItems,
      modalDialogRef,
      map,
      date,
      onUpdateTable,
      onDeleteItem,
      doRefresh,
      onUpdateBorder,
      onUpdateStriped,
      rowClassNameFun,
      onUpdateItem,
      onDataFormConfirm,
      check,
    }
  },
})
</script>

<style lang="less" scoped>
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
<style scoped>
.ant-table-striped :deep(.table-striped) td {
  background-color: #fafafa;
}
</style>
