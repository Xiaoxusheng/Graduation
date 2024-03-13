<template>
  <div class="main-container">
    <TableBody>
      <template #header>
        <TableHeader :show-filter="false">
          <template #table-config>
            <AddButton @add="onAddItem"/>
          </template>
        </TableHeader>
      </template>
      <template #default>
        <a-table
            :bordered="{ wrapper: true, cell: true }"
            :data="dataList"
            :loading="tableLoading"
            :pagination="false"
            :row-key="rowKey"
        >
          <template #columns>
            <a-table-column
                v-for="item of tableColumns"
                :key="item.key"
                :align="item.align"
                :data-index="(item.key as string)"
                :fixed="item.fixed"
                :title="(item.title as string)"
                :width="item.width"
            >
              <template v-if="item.key === 'index'" #cell="{ rowIndex }">
                {{ rowIndex + 1 }}
              </template>
              <template v-if="item.key === 'name'" #cell="{ record }">
                {{ record }}
              </template>
              <template v-if="item.key === 'actions'" #cell="{ record }">
                <a-space>
                  <a-button size="mini" status="success" @click="onUpdateItem(record)"
                  >编辑
                  </a-button>
                  <a-button size="mini" status="danger" @click="onDeleteItem(record)"
                  >删除
                  </a-button>
                  <a-button size="mini" status="warning" @click="onShowMenu(record)">
                    菜单权限
                  </a-button>
                </a-space>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </template>
    </TableBody>
    <ModalDialog ref="modalDialogRef" :title="actionTitle" @confirm="onDataFormConfirms">
      <template #content>
        <a-form :model="formModel">
          <a-form-item
              v-for="item of formItems"
              :key="item.key"
              :class="[item.required ? 'form-item__require' : 'form-item__no_require']"
              :label="item.label"
          >
            <template v-if="item.type === 'input'">
              <a-input v-model="item.value.value" :placeholder="item.placeholder">
                <template v-if="item.key === 'roleCode'" #prepend>
                  {{ 123 }}
                </template>
              </a-input>
            </template>
          </a-form-item>
        </a-form>
      </template>
    </ModalDialog>
    <ModalDialog ref="menuModalDialogRef" title="编辑菜单权限" @confirm="onDataFormConfirm(menuData)">
      <template #content>
        <a-tree
            ref="tree"
            v-model:checked-keys="defaultCheckedKeys"
            v-model:expanded-keys="defaultExpandedKeys"
            @check="shows "
            :data="menuData"
            checkable
        />
      </template>
    </ModalDialog>
  </div>
</template>

<script lang="ts">
import {get, post} from '@/api/http'
import {addRolesForUser, deleteRole, getRoleMenuList, roleList, updateRoleMenu} from '@/api/url'
import {useRowKey, useTable, useTableColumn} from '@/hooks/table'
import {FormItem, ModalDialogType} from '@/types/components'
import {Message, Modal,} from '@arco-design/web-vue'
import {defineComponent, nextTick, onMounted, ref} from 'vue'
import useUserStore from "@/store/modules/user";

const ROLE_CODE_FLAG = 'ROLE_'
const formItems = [
  {
    label: '角色名称',
    type: 'input',
    key: 'name',
    value: ref(''),

    required: true,
    placeholder: '请输入角色名称',
    validator: function () {
      if (!this.value.value) {
        Message.error(this.placeholder || '')
        return false
      }
      return true
    },
  },

  {
    label: '角色id',
    key: 'description',
    value: ref(''),
    type: 'input',
    placeholder: '请输入id',
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
const tree = ref<null>(null)

function handleMenuData(
    menuData: Array<any>,
    defaultCheckedKeys: Array<string>,
    defaultExpandedKeys: Array<string>
) {
  const tempMenus = [] as Array<any>
  menuData.forEach((it) => {
    const tempMenu = {} as any
    tempMenu.key = it.menuUrl
    tempMenu.title = it.menuName
    tempMenu.menu = it.uid
    tempMenu.ok = it.is

    if (it.is) {
      defaultCheckedKeys.push(tempMenu.key as string)
    }
    if (it.children) {
      defaultExpandedKeys.push(tempMenu.key as string)
      tempMenu.children = handleMenuData(it.children, defaultCheckedKeys, defaultExpandedKeys)
    }
    tempMenus.push(tempMenu)
  })
  return tempMenus
}

export default defineComponent({
  name: 'Role',
  setup() {
    const modalDialogRef = ref<ModalDialogType | null>(null)
    const menuModalDialogRef = ref<ModalDialogType | null>(null)
    const table = useTable()
    const rowKey = useRowKey('id')
    const actionTitle = ref('添加角色')
    const menuData = ref([] as Array<any>)
    const tableColumns = useTableColumn([
      table.indexColumn,
      {
        title: '角色名称',
        key: 'name',
        dataIndex: 'name',
        width: 300
      },
      {
        title: '操作',
        key: 'actions',
        dataIndex: 'actions',
        width: 300

      },
    ])
    const defaultCheckedKeys = ref([] as Array<string>)
    const defaultExpandedKeys = ref([] as Array<string>)
    const userStore = useUserStore()
    const formModel = ref({})
    const res = [] as any
    let role = ""

    function doRefresh() {
      get({
        url: roleList,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
      })
          .then(table.handleSuccess)
          .catch(console.log)
    }

    function onAddItem() {
      actionTitle.value = '添加角色'
      modalDialogRef.value?.toggle()
      formItems.forEach((it) => {
        if (it.reset) {
          it.reset()
        } else {
          it.value.value = ''
        }
      })
    }

    function onUpdateItem(item: any) {
      actionTitle.value = '编辑角色'
      modalDialogRef.value?.toggle()
      nextTick(() => {
        formItems.forEach((it) => {
          const key = it.key
          const propName = item[key]
          if (propName) {
            if (it.key === 'roleCode') {
              it.value.value = propName.replace(ROLE_CODE_FLAG, '')
            } else {
              it.value.value = propName
            }
          }
        })
      })
    }

    function onDeleteItem(data: any) {
      Modal.confirm({
        title: '提示',
        content: '是否要删除此角色？',
        cancelText: '取消',
        okText: '删除',
        onOk: () => {
          //   修改菜单
          post({
            url: deleteRole,
            headers: {
              Authorization: "Bearer " + userStore.token,
              'Content-Type': "application/x-www-form-urlencoded; charset=UTF-8"
            },
            data: {
              role: data,
            }
          }).then((res) => {
                doRefresh()
                Message.success("删除成功！")
              }
          ).catch(error => {
            Message.error(error.message)
          })
        },
      })
    }

    function onDataFormConfirm(item: any) {
      menuModalDialogRef.value?.toggle()
      const data = tree.value.getCheckedNodes()
      res.length = 0
      console.log(data)
      data.forEach((i: any) => {
        if (i.ok) {
          const resp = i.key.split("/")
          menuData.value.forEach((j => {
            if (j.key == "/" + resp[1] && resp.length > 2) {
              if (!res.includes(j.menu.toString())) {
                res.push(j.menu.toString())
              }
            }
          }))
          res.push(i.menu.toString())
        }
      })

      //   修改菜单
      post({
        url: updateRoleMenu,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: {
          role: role,
          menu: [...new Set(res)],
        }
      }).then((res) => {
            Message.success("更新成功！")
          }
      ).catch(error => {
        Message.error(error.message)
      })
    }

    function onDataFormConfirms(item: any) {
      modalDialogRef.value?.toggle()

      //   修改菜单
      post({
        url: addRolesForUser,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: {
          role: formItems[0].value.value,
          user: formItems[1].value.value,
        }
      }).then((res) => {
            Message.success("更新成功！")
          }
      ).catch(error => {
        Message.error(error.message)
      })
      console.log()
    }

    function onShowMenu(item: any) {
      role = item
      get({
        url: getRoleMenuList,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: {
          role: item,
        },
      })
          .then((res) => {
            menuData.value = []
            defaultCheckedKeys.value = []
            defaultExpandedKeys.value = []
            menuData.value = handleMenuData(
                res.data,
                defaultCheckedKeys.value,
                defaultExpandedKeys.value
            )
            menuModalDialogRef.value?.toggle()
          })
          .catch(console.log)
    }

    function shows(checkedKeys: any, data: any) {
      data.node.ok = !data.node.ok
      if (data.node.children) {
        console.log(data.node.children.length)
        data.node.children.forEach((i: any) => {
          console.log(i)
          i.ok = !i.ok
        })
      }
      console.log(data)
    }

    onMounted(doRefresh)
    return {
      tree,
      ROLE_CODE_FLAG,
      modalDialogRef,
      menuModalDialogRef,
      rowKey,
      formModel,
      menuData,
      tableColumns,
      formItems,
      actionTitle,
      defaultCheckedKeys,
      defaultExpandedKeys,
      ...table,
      role,
      shows,
      onAddItem,
      onDataFormConfirm,
      onDataFormConfirms,
      onShowMenu,
      onDeleteItem,
      onUpdateItem,
    }
  },
})
</script>
